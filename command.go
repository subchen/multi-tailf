package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/ttacon/chalk"
)

type (
	Command struct {
		prompts string
		cmd     *exec.Cmd
	}
)

func (c *Command) Execute(output chan string) {
	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := c.cmd.Start(); err != nil {
		panic(fmt.Sprintf("[%s] failed to start: %s", c.prompts, err))
	}

	c.watch(stdout, output, false)
	c.watch(stderr, output, true)

	if err := c.cmd.Wait(); err != nil {
		fmt.Println(c.prompts + strings.Join(c.cmd.Args, " "))
		panic(err)
	}
}

func (c *Command) watch(input io.ReadCloser, output chan string, isErr bool) {
	go func() {
		reader := bufio.NewReader(input)
		for {
			line, err := reader.ReadString('\n')
			if err != nil || io.EOF == err {
				if err != io.EOF {
					panic(fmt.Sprintf("[%s] read error: %s", c.prompts, err))
				}
				break
			}

			if isErr {
				line = chalk.Red.Color(line)
			}

			output <- c.prompts + line
		}
	}()
}
