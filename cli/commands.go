package cli

import (
	"errors"
	"fmt"
)

type Runnable interface {
	Run(args []string) error
}

type RunFunc func(args []string) error

func (r RunFunc) Run(args []string) error {
	return r(args)
}

type CmdRegistry map[string]Runnable

func (c CmdRegistry) Register(name string, cmd Runnable) {
	if _, ok := c[name]; ok {
		panic(fmt.Errorf("subcommand %s already registered", name))
	}
	c[name] = cmd
}

var ErrCommandNotFound = errors.New("no subcommand registered with that name")
var ErrNoSubcommandSupplied = errors.New("no subcommand supplied")
var ErrUnimplemented = errors.New("subcommand is unimplemented")

func (c CmdRegistry) Run(args []string) error {
	if len(args) < 1 {
		return ErrNoSubcommandSupplied
	}
	cn, args := args[0], args[1:]
	cmd, ok := c[cn]
	if !ok {
		return ErrCommandNotFound
	}
	return cmd.Run(args)
}

func (c CmdRegistry) Usage() {
	fmt.Println("Subcommands: ")
	for key, _ := range c {
		fmt.Printf("%s\n", key)
	}
}

func Unimplemented(args []string) error {
	return ErrUnimplemented
}
