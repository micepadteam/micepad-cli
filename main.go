package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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
		case "version":
			handleVersion()
			return
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
}

func handleVersion() {
	fmt.Printf("micepad %s (%s)\n", version, commit)
	fmt.Printf("Server:  %s\n", config.ResolveURL())
	fmt.Printf("Config:  %s\n", config.Path())
	fmt.Printf("Storage: %s\n", config.Dir())
}

const (
	repo          = "micepadteam/micepad-cli"
	installScript = "https://github.com/" + repo + "/releases/latest/download/install.sh"
)

func getLatestVersion() (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get("https://github.com/" + repo + "/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to check latest version: %w", err)
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("could not determine latest version (no redirect)")
	}

	re := regexp.MustCompile(`/v?(\d+\.\d+\.\d+.*)$`)
	matches := re.FindStringSubmatch(location)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not parse version from %s", location)
	}
	return matches[1], nil
}

func handleUpdate() {
	fmt.Printf("Current version: %s (%s)\n", version, commit)
	fmt.Println("Checking for updates...")

	latest, err := getLatestVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not check latest version: %v\n", err)
		fmt.Fprintln(os.Stderr, "Proceeding with update anyway...")
	} else if latest == version {
		fmt.Printf("Already up to date (v%s)\n", version)
		updateSkill()
		return
	} else {
		fmt.Printf("New version available: v%s → v%s\n", version, latest)
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL %s | bash", installScript))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		os.Exit(1)
	}
}

func updateSkill() {
	if _, err := exec.LookPath("npx"); err != nil {
		return
	}

	// Check if skill updates are available
	check := exec.Command("npx", "-y", "skills", "check")
	output, err := check.CombinedOutput()
	if err != nil {
		return
	}

	if !strings.Contains(string(output), "update(s) available") {
		fmt.Println("Skills already up to date.")
		return
	}

	fmt.Println("Updating Micepad skills for Claude Code...")
	cmd := exec.Command("npx", "-y", "skills", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Skill update skipped (non-critical)")
	}
}
