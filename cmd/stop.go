package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type Stop struct{}

func (c *Stop) Run(args []string) error {
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

	if _, err := fmt.Fprintf(f, "%d ", time.Now().Unix()); err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	//TODO default
	tmpfile, err := ioutil.TempFile("", "clocking")
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	defer os.Remove(tmpfile.Name())

	edt := exec.Command(os.Getenv("EDITOR"), tmpfile.Name())
	edt.Stdin = os.Stdin
	edt.Stdout = os.Stdout
	edt.Stderr = os.Stderr

	if err := edt.Run(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	msg, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}

	if _, err := fmt.Fprintf(f, "%q", msg); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	if err := save(string(msg)); err != nil {
		log.Println("[DEBUG]", err)
	}

	return nil
}

func save(msg string) error {
	g := &Git{}
	if err := g.Run([]string{
		"add",
		".",
	}); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	return g.Run([]string{
		"commit",
		"-a",
		"-m",
		msg,
	})
}
