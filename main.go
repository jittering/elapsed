package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
)

var defaultDateTimeLayout = time.RFC3339
var defaultSlowThreshold int64 = 500

var showTS = ""
var showElapsed = true
var showDelta = true
var slow int64 = 0

// getArgs from CLI or config file
func getArgs() []string {
	args := os.Args[1:]
	if len(args) != 0 {
		return args
	}

	// no cli args given, look for config file
	dir, err := homedir.Dir()
	if err != nil {
		return args
	}
	rc := filepath.Join(dir, ".elapsedrc")
	_, err = os.Stat(rc)
	if err != nil {
		return args
	}
	b, err := os.ReadFile(rc)
	if err != nil {
		return args
	}
	return strings.Split(strings.TrimSpace(string(b)), " ")
}

func parseFlags() {
	st := flag.Bool("datetime", false, "Show date/time stamp when message was received")
	stf := flag.String("format", "", "Date/time format (default: "+defaultDateTimeLayout+")")

	se := flag.Bool("no-elapsed", false, "Do not print the absolute elapsed time")
	sd := flag.Bool("no-delta", false, "Do not print the delta elapsed time")

	sl := flag.Bool("slow", false, "Show only slow deltas (over a certain threshold)")
	slm := flag.Int64("slow-ms", defaultSlowThreshold, "Slow delta threshold in ms")

	flag.CommandLine.Parse(getArgs())

	if se != nil && *se {
		showElapsed = false
	}
	if sd != nil && *sd {
		showDelta = false
	}
	if st != nil && *st {
		if stf != nil && *stf != "" {
			showTS = *stf
		} else {
			showTS = defaultDateTimeLayout
		}
	}
	if sl != nil && *sl {
		slow = defaultSlowThreshold
	}
	if slm != nil && *slm != defaultSlowThreshold {
		slow = *slm
	}
}

func run(r io.Reader) {
	f := "["
	if showTS != "" {
		f += "%s"
		if showElapsed || showDelta {
			f += " | "
		}
	}
	if showElapsed {
		f += "%-7s"
	}
	if showDelta {
		if showElapsed {
			f += " "
		}
		f += "%-8s"
	}
	f += "] %s"

	reader := bufio.NewReader(r)
	start := time.Now()
	last := time.Now()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		var args []interface{}

		if showTS != "" {
			// add current date/time stamp
			args = append(args, time.Now().Format(showTS))
		}

		if showElapsed {
			// absolute elapsed time
			elapsed := time.Since(start).Round(1 * time.Microsecond).Milliseconds()
			elapsedS := fmt.Sprintf("%dms", elapsed)
			args = append(args, elapsedS)
		}

		if showDelta {
			// delta between last line
			delta := time.Since(last).Round(1 * time.Microsecond).Milliseconds()
			if delta >= slow {
				deltaS := fmt.Sprintf("+%dms", delta)
				args = append(args, deltaS)
			} else {
				args = append(args, "")
			}
		}

		args = append(args, line)
		// "[%-7s +%-7s] %s"
		fmt.Printf(f, args...)

		if showDelta {
			// update last line ts
			last = time.Now()
		}
	}
}

func getReader() (io.Reader, []string) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// look for cmd
		args := os.Args[1:]
		if len(args) != 0 {
			for i, arg := range args {
				if arg == "--" && i+1 < len(args) && args[i+1] != "" {
					return nil, args[i+1:]
				}
			}
		}

		fmt.Println("error: stdin not connected to a pipe")
		os.Exit(1)
	}

	return os.Stdin, nil
}

func runCmd(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := cmd.Start()
		if err != nil {
			fmt.Println("error running command: ", err)
		}
		go run(r)
		err = cmd.Wait()
		if err != nil {
			fmt.Println("error waiting for command: ", err)
		}
		os.Exit(cmd.ProcessState.ExitCode())
	}()
	wg.Wait()
}

func main() {
	parseFlags()
	if !(showElapsed || showDelta) {
		fmt.Println("error: must enable either elapsed or delta")
		os.Exit(1)
	}

	r, args := getReader()
	if args != nil {
		runCmd(args)
	} else {
		run(r)
	}
}
