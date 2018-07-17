package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/rjeczalik/notify"
)

func main() {
	//procs, _ := ps.Processes()
	//for _, proc := range procs {
	//        log.Println(proc.Executable())
	//}
	dir := "/home/whytheplatypus"
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return nil
		}
		if info.IsDir() {
			count = count + 1
		}
		return nil
	})
	log.Println(count)

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}
	fmt.Println("vim-go")
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	go func() {
		<-s
		log.Println("Interupt received, shutting down...")
		close(c)
	}()

	// Set up a watchpoint listening for inotify-specific events within a
	// current working directory. Dispatch each InCloseWrite and InMovedTo
	// events separately to c.
	if err := notify.Watch(
		"/home/whytheplatypus/Development/fearless/bluebutton/bluebutton-web-server",
		c,
		//notify.InOpen,
		notify.InAccess,
		//notify.InCloseNowrite,
		notify.Write,
		//notify.InCloseWrite,
	); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	// Block until an event is received.
	for ei := range c {
		log.Printf("%+v \n", ei)
	}
}
