package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func main() {
	patterns := os.Args[1:]
	procs := map[string]string{}
	for {
		ps, err := Getprocs()
		if err != nil {
			log.Fatal(err)
		}
		remaining := map[string]struct{}{}
		for _, p := range ps {
			_, ok := procs[p]
			if !ok {
				path, err := Getcwd(p)
				if err != nil {
					continue
				}
				if !match(patterns, path) {
					continue
				}
				exe, _ := Getexe(p)
				log.Println("found ", p, exe, path)
				procs[p] = path
			}
			remaining[p] = struct{}{}
		}

		for p, path := range procs {
			_, ok := remaining[p]
			if !ok {
				log.Println("lost ", path)
				delete(procs, p)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func match(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if ok, _ := regexp.MatchString(pattern, path); ok {
			return true
		}
	}
	return false
}

func Getcwd(pid string) (s string, err error) {
	return os.Readlink(
		fmt.Sprintf("%s/cwd", pid))
}

func Getexe(pid string) (s string, err error) {
	return os.Readlink(
		fmt.Sprintf("%s/exe", pid))
}

func Getprocs() ([]string, error) {
	return filepath.Glob("/proc/[0-9]*")
}
