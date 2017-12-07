package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/whytheplatypus/clocking/timesheet"
)

type Backfill struct{}

func (c *Backfill) Run(args []string) error {
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
	tmpfile, err := ioutil.TempFile("", "clocking")
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	defer os.Remove(tmpfile.Name())
	p := &timesheet.Punch{}
	for scanner.Scan() {
		if len(scanner.Bytes()) > 0 {
			if err := timesheet.UnmarshalCLK(scanner.Bytes(), p); err != nil {
				log.Println("[ERROR]", err)
				return err
			}
			t, err := p.MarshalTime(time.RFC3339)
			if err != nil {
				log.Println("[ERROR]", err)
				return err
			}
			if _, err := fmt.Fprintln(tmpfile, string(t)); err != nil {
				log.Println("[ERROR]", err)
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	f.Close()

	edt := exec.Command(os.Getenv("EDITOR"), tmpfile.Name())
	edt.Stdin = os.Stdin
	edt.Stdout = os.Stdout
	edt.Stderr = os.Stderr

	if err := edt.Run(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	tc, err := timesheet.ReadFile(tmpfile.Name(), timesheet.UnmarshalTime(time.RFC3339))
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	sort.Sort(tc)

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
	// commit
	return nil
}
