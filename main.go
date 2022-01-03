package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("error: stdin not connected to a pipe")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	start := time.Now()
	last := time.Now()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		elapsed := time.Since(start).Round(1 * time.Microsecond).Milliseconds()
		elapsedS := fmt.Sprintf("%dms", elapsed)
		delta := time.Since(last).Round(1 * time.Microsecond).Milliseconds()
		deltaS := fmt.Sprintf("%dms", delta)
		fmt.Printf("[%-6s +%-6s] %s", elapsedS, deltaS, line)
		last = time.Now()
	}
}
