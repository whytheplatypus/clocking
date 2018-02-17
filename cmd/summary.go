package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/whytheplatypus/clocking/timesheet"
)

type Summary struct{}

func (c *Summary) Run(args []string) error {
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
	var m time.Month
	days := map[string]struct {
		m string
		t float64
	}{}
	prev := struct {
		m string
		t float64
	}{}
	var prevd string
	for scanner.Scan() {
		if len(scanner.Bytes()) > 0 {
			if err := timesheet.UnmarshalCLK(scanner.Bytes(), p); err != nil {
				log.Println("[ERROR]", err)
				return err
			}
			s := time.Unix(p.Start, 0)
			// bug where year is different but month is the same
			// not risky because of ordering but still
			t := time.Unix(p.End, 0).Sub(time.Unix(p.Start, 0))

			//if days[s.Format("01/02/2006")] == nil {
			//        days[s.Format("01/02/2006")] = struct {
			//                m string
			//                t int
			//        }{}
			//}
			if days[s.Format("01/02/2006")] == struct {
				m string
				t float64
			}{} {
				if prevd != "" {
					fmt.Println(prevd, prev.m, prev.t)
				}
			}
			d := days[s.Format("01/02/2006")]
			d.m += strings.Replace(p.Msg, "\\n", " ", -1)
			d.t += t.Round(time.Duration(time.Hour / 4)).Hours()

			days[s.Format("01/02/2006")] = d
			prev = d
			prevd = s.Format("01/02/2006")
			//fmt.Println(
			//s.Format("01/02/2006"),
			//strings.Replace(p.Msg, "\\n", " ", -1),
			//t.Round(time.Duration(time.Hour/4)).Hours(),
			//)
			if s.Month() != m {
				m = s.Month()
				fmt.Println("\n\n-----", s.Format("01/2006"), "------")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	return nil
}
