package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/whytheplatypus/clocking/cli"
	"github.com/whytheplatypus/clocking/cmd"
)

var (
	usage = `
clocking usage
`
	subcmds = cli.CmdRegistry{}
)

func main() {
	fmt.Println("The clock king")
	subcmds.Register(
		"in",
		&cmd.Start{})
	subcmds.Register(
		"out",
		&cmd.Stop{})
	subcmds.Register(
		"git",
		&cmd.Git{})
	subcmds.Register(
		"backfill",
		&cmd.Backfill{})
	subcmds.Register(
		"calculate",
		&cmd.Calculate{})
	subcmds.Register(
		"summary",
		&cmd.Summary{})

	var verbose bool
	cmdFlag := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cmdFlag.BoolVar(&verbose, "v", false, "Enable for verbose logging")
	if err := cmdFlag.Parse(os.Args[1:]); err != nil {
		log.Println("[ERROR]", err)
		subcmds.Usage()
		os.Exit(1)
	}

	if verbose {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	args := cmdFlag.Args()
	log.Println("[DEBUG]", os.Args, args)

	if err := subcmds.Run(args); err != nil {
		log.Println("[ERROR]", err)
		subcmds.Usage()
		os.Exit(1)
	}

	os.Exit(0)

}
