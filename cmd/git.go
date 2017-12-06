package cmd

import (
	"log"
	"os"
	"os/exec"
)

type Git struct{}

func (c *Git) Run(args []string) error {
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

	edt := exec.Command("git", args...)
	edt.Stdin = os.Stdin
	edt.Stdout = os.Stdout
	edt.Stderr = os.Stderr

	if err := edt.Run(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	return nil
}
