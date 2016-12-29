package main

import (
	"fmt"
	"sync"
)

type (
	CommandList struct {
		output   chan string
		wg       sync.WaitGroup
		commands []Command
	}
)

func (cl *CommandList) AppendCommand(c *Command) {
	cl.commands = append(cl.commands, *c)
}

func (cl *CommandList) Start() {
	cl.output = make(chan string, 50)

	cl.watch()

	for _, c := range cl.commands {
		cl.execCommand(c)
	}
	cl.wg.Wait()
}

func (cl *CommandList) watch() {
	go func() {
		for line := range cl.output {
			fmt.Print(line)
		}
	}()
}

func (cl *CommandList) execCommand(c Command) {
	cl.wg.Add(1)
	go func() {
		c.Execute(cl.output)
		defer cl.wg.Done()
	}()
}
