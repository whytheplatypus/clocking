package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
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
		if err := p.UnmarshalCLK(scanner.Bytes()); err != nil {
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
	log.Print("[DEBUG]", "tmp file: ")
	l, _ := ioutil.ReadFile(tmpfile.Name())
	log.Println(string(l))

	// read and parse result of tmp file
	tscanner := bufio.NewScanner(bytes.NewBuffer(l))
	tc := timesheet.TimeCard{}
	for tscanner.Scan() {
		log.Println("[DEBUG]", string(tscanner.Bytes()))
		if len(tscanner.Bytes()) < 1 {
			continue
		}
		p := &timesheet.Punch{}
		log.Println("[DEBUG]", string(tscanner.Bytes()))
		if err := p.UnmarshalTime(tscanner.Bytes(), time.RFC3339); err != nil {
			log.Println("[ERROR]", err)
			return err
		}
		tc = append(tc, p)
	}
	if err := tscanner.Err(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	sort.Sort(tc)
	// replace contents of project file
	//TODO put in tmp file intermediary so it's transactional
	ff, err := os.OpenFile(project, os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	defer ff.Close()
	t := template.Must(template.New("timesheet").Parse(timesheet.Template))
	if err := t.Execute(ff, tc); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	// commit
	return nil
}
