package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/micepadteam/micepad-cli/internal/config"
	"github.com/micepadteam/micepad-cli/internal/terminalwire"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Handle local flags (before connecting to server)
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("micepad %s (%s)\n", version, commit)
		return
	}

	// Handle local commands (before connecting to server)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "configure":
			handleConfigure(os.Args[2:])
			return
		case "update":
			handleUpdate()
			return
		}
	}

	wsURL := config.ResolveURL()

	client, err := terminalwire.Connect(wsURL, "micepad", version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Server: %s\n", wsURL)
		fmt.Fprintf(os.Stderr, "Run 'micepad configure' to change the server URL.\n")
		os.Exit(1)
	}

	if err := client.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleConfigure(args []string) {
	cfg := config.Load()

	// Handle --url flag for non-interactive URL setting
	for i, arg := range args {
		if arg == "--url" && i+1 < len(args) {
			cfg.URL = args[i+1]
			if err := cfg.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("URL set to %s\n", cfg.URL)
			fmt.Printf("Configuration saved to %s\n", config.Path())
			return
		}
	}

	// Interactive configuration
	currentURL := config.ResolveURL()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Server URL [%s]: ", currentURL)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		cfg.URL = input
	} else {
		cfg.URL = currentURL
	}

	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration saved to %s\n", config.Path())

	// Connect to server for account/event selection
	wsURL := config.ResolveURL()
	client, err := terminalwire.Connect(wsURL, "micepad", version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCould not connect to server: %v\n", err)
		fmt.Println("URL saved. Run 'micepad configure' again once the server is available.")
		return
	}

	if err := client.Run([]string{"configure"}); err != nil {
		// Server may not support 'configure' command yet — URL is already saved
	}
}

const installScript = "https://raw.githubusercontent.com/micepadteam/micepad-cli/master/scripts/install.sh"

func handleUpdate() {
	fmt.Printf("Current version: %s (%s)\n", version, commit)
	fmt.Println("Checking for updates...")

	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL %s | bash", installScript))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}
}
