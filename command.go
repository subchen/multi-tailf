package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ttacon/chalk"
)

type (
	Command struct {
		prompts           string
		command           string
		commandParameters []string
		wg                sync.WaitGroup
		retry             bool
	}
)

func (c *Command) Execute(output chan string) {
	for {
		c.retry = false

		output <- c.prompts + chalk.Blue.Color(fmt.Sprintf("%s %s\n", c.command, strings.Join(c.commandParameters, " ")))

		cmd := exec.Command(c.command, c.commandParameters...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}

		if err := cmd.Start(); err != nil {
			panic(fmt.Sprintf("[%s] failed to start: %s", c.prompts, err))
		}

		c.watch(stdout, output, false)
		c.watch(stderr, output, true)

		if err := cmd.Wait(); err != nil {
			c.wg.Wait()
			if !c.retry {
				//panic(err)
				output <- c.prompts + err.Error() + "\n"
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func (c *Command) watch(input io.ReadCloser, output chan string, isErr bool) {
	c.wg.Add(1)
	go func() {
		reader := bufio.NewReader(input)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				//if err != io.EOF {
				//	panic(fmt.Sprintf("[%s] read error: %s", c.prompts, err))
				//}
				break
			}

			if isErr {
				if strings.HasPrefix(line, "tail: cannot open") && strings.HasSuffix(line, "for reading: No such file or directory\n") {
					c.retry = true
				}
				line = chalk.Red.Color(line)
			}

			output <- c.prompts + line
		}

		c.wg.Done()
	}()
}
