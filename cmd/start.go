package cmd

import (
	"fmt"
	"log"
	"os"
	"time"
)

const projdir = ".clocking"

type Start struct{}

func (c *Start) Run(args []string) error {
	// if it doesn't exist, make a directory to store records in
	home, ok := os.LookupEnv("HOME")
	if !ok {
		home = "."
	}
	os.Chdir(home)
	if err := os.Mkdir(projdir, 0700); err != nil && err != os.ErrExist {
		log.Println("[DEBUG]", err)
	}
	if err := os.Chdir(projdir); err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	// get the project name
	log.Println("[DEBUG]", args)
	// TODO handle this doesn't exist
	if len(args) < 1 {
		return fmt.Errorf("No project specified")
	}

	project := args[0]
	// open the project file
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(project, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%d:", time.Now().Unix()); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	// make a new entry
	// error if there's already an active entry
	return nil
}
