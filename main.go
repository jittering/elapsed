package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

var defaultDateTimeLayout = time.RFC3339
var defaultSlowThreshold int64 = 500

var showTS = ""
var showElapsed = true
var showDelta = true
var slow int64 = 0

func parseFlags() {
	st := flag.Bool("datetime", false, "Show date/time stamp when message was received")
	stf := flag.String("format", "", "Date/time format (default: "+defaultDateTimeLayout+")")

	se := flag.Bool("no-elapsed", false, "Do not print the absolute elapsed time")
	sd := flag.Bool("no-delta", false, "Do not print the delta elapsed time")

	sl := flag.Bool("slow", false, "Show only slow deltas (over a certain threshold)")
	slm := flag.Int64("slow-ms", defaultSlowThreshold, "Slow delta threshold in ms")

	flag.Parse()

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

func main() {
	parseFlags()
	if !(showElapsed || showDelta) {
		fmt.Println("error: must enable either elapsed or delta")
		os.Exit(1)
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("error: stdin not connected to a pipe")
		os.Exit(1)
	}

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

	reader := bufio.NewReader(os.Stdin)
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
