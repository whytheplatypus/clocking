package timesheet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"text/template"
	"time"
)

var t *template.Template

func init() {
	t = template.Must(template.New("timesheet").Parse(Template))
}

type TimeCard []*Punch

func (a TimeCard) Len() int           { return len(a) }
func (a TimeCard) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TimeCard) Less(i, j int) bool { return a[i].Start < a[j].Start }

func (tc TimeCard) Execute(w io.Writer) error {
	return t.Execute(w, tc)
}

type Parser func([]byte, *Punch) error

func ReadFile(f string, prsr Parser) (TimeCard, error) {
	tc := TimeCard{}
	log.Print("[DEBUG]", "tmp file: ")
	l, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println("[ERROR]", err)
		return tc, err
	}
	log.Println(string(l))

	// read and parse result of tmp file
	tscanner := bufio.NewScanner(bytes.NewBuffer(l))
	for tscanner.Scan() {
		log.Println("[DEBUG]", string(tscanner.Bytes()))
		if len(tscanner.Bytes()) < 1 {
			continue
		}
		p := &Punch{}
		log.Println("[DEBUG]", string(tscanner.Bytes()))
		if err := prsr(tscanner.Bytes(), p); err != nil {
			log.Println("[ERROR]", err)
			return tc, err
		}
		tc = append(tc, p)
	}
	if err := tscanner.Err(); err != nil {
		log.Println("[ERROR]", err)
		return tc, err
	}
	return tc, nil
}

type Punch struct {
	Start int64
	End   int64
	Msg   string
}

var ErrBadFmt = errors.New("There was a problem with the record time format")

func UnmarshalCLK(in []byte, p *Punch) error {
	log.Println("[DEBUG]", string(in))
	times := bytes.SplitN(in, []byte{':'}, 2)
	start := times[0]
	stpmsg := bytes.SplitN(times[1], []byte{' '}, 2)
	stop := stpmsg[0]

	msg := stpmsg[1]

	log.Println("[DEBUG]", "start string:", start)
	var err error
	p.Start, err = strconv.ParseInt(string(start), 10, 64)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	log.Println("[DEBUG]", "start:", p.Start)
	p.End, err = strconv.ParseInt(string(stop), 10, 64)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	log.Println("[DEBUG]", "stop:", p.End)

	p.Msg = string(msg)

	return nil
}

func (p *Punch) MarshalText() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	if _, err := fmt.Fprintf(b, "%d:%d %s", p.Start, p.End, p.Msg); err != nil {
		log.Println("[ERROR]", err)
		return b.Bytes(), err
	}
	return b.Bytes(), nil
}

func UnmarshalTime(layout string) Parser {
	return func(in []byte, p *Punch) error {
		parts := bytes.Split(in, []byte{'|'})
		start, err := time.Parse(layout, string(bytes.TrimSpace(parts[0])))
		if err != nil {
			log.Println("[ERROR]", err)
			return err
		}
		stop, err := time.Parse(layout, string(bytes.TrimSpace(parts[1])))
		if err != nil {
			log.Println("[ERROR]", err)
			return err
		}
		p.Start = start.Unix()
		p.End = stop.Unix()
		p.Msg = string(bytes.TrimSpace(parts[2]))
		return nil
	}
}

func (p *Punch) MarshalTime(layout string) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	strt := time.Unix(p.Start, 0).Local()
	strts := strt.Format(layout)
	stp := time.Unix(p.End, 0).Local()
	stps := stp.Format(layout)
	if _, err := fmt.Fprintf(b, "%s | %s | %s", strts, stps, p.Msg); err != nil {
		log.Println("[ERROR]", err)
		return b.Bytes(), err
	}
	return b.Bytes(), nil
}
