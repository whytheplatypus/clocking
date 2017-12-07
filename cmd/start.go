package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/whytheplatypus/clocking/timesheet"
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

	tc, err := timesheet.ReadFile(project, timesheet.UnmarshalCLK)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	tc = append(tc, &timesheet.Punch{
		Start: time.Now().Unix(),
	})

	ff, err := os.OpenFile(project, os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	defer ff.Close()

	if err := tc.Execute(ff); err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	return nil
}
