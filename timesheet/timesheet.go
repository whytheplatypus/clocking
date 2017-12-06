package timesheet

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type TimeCard []*Punch

func (a TimeCard) Len() int           { return len(a) }
func (a TimeCard) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TimeCard) Less(i, j int) bool { return a[i].Start < a[j].Start }

type Punch struct {
	Start int64
	End   int64
	Msg   string
}

var ErrBadFmt = errors.New("There was a problem with the record time format")

func (p *Punch) UnmarshalCLK(in []byte) error {
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

func (p *Punch) UnmarshalTime(in []byte, layout string) error {
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
