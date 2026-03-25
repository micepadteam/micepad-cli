package main

import (
	"fmt"
	"os"

	"github.com/micepadteam/micepad-cli/internal/terminalwire"
)

var (
	version = "dev"
	commit  = "none"
)

const defaultURL = "wss://studio.micepad.co/terminal"

func main() {
	// Handle local version flag (before connecting to server)
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("micepad %s (%s)\n", version, commit)
		return
	}

	wsURL := defaultURL
	if envURL := os.Getenv("MICEPAD_URL"); envURL != "" {
		wsURL = envURL
	}

	client, err := terminalwire.Connect(wsURL, "micepad")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		os.Exit(1)
	}

	if err := client.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
