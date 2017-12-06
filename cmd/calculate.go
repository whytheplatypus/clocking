package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/whytheplatypus/clocking/timesheet"
)

type Calculate struct{}

func (c *Calculate) Run(args []string) error {
	// open and read the project file
	home, ok := os.LookupEnv("HOME")
	if !ok {
		home = "."
	}
	os.Chdir(home)
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
	f, err := os.Open(project)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// write tmp file
	p := &timesheet.Punch{}
	t := time.Duration(0)
	for scanner.Scan() {
		if err := p.UnmarshalCLK(scanner.Bytes()); err != nil {
			log.Println("[ERROR]", err)
			return err
		}
		t += time.Unix(p.End, 0).Sub(time.Unix(p.Start, 0))
	}
	if err := scanner.Err(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	fmt.Println(t)
	return nil
}
