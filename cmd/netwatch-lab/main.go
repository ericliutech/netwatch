package main

import (
	"fmt"
	"os"
)

type runner func() error

var runners = map[string]runner{
	"arp":       runARP,
	"tcp-probe": runTCPProbe,
	"udp-probe": runUDPProbe,
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	name := os.Args[1]

	run, ok := runners[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown runner: %s\n\n", name)
		printUsage()
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("usage: go run ./cmd/netwatch-lab <runner>")
	fmt.Println()
	fmt.Println("available runners:")

	for name := range runners {
		fmt.Printf("  %s\n", name)
	}
}
