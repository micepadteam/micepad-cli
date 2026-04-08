package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/micepad/micepad-cli/internal/config"
	"github.com/micepad/micepad-cli/internal/terminalwire"
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

	// Extract -e/--env flag from anywhere in the args.
	envOverride, args := extractEnvFlag(os.Args[1:])

	// Handle local commands (before connecting to server)
	if len(args) > 0 {
		switch args[0] {
		case "version":
			handleVersion()
			return
		case "configure":
			handleConfigure(args[1:])
			return
		case "update":
			handleUpdate()
			return
		case "env":
			handleEnv(args[1:])
			return
		}
	}

	var wsURL string
	if envOverride != "" {
		wsURL = config.ResolveURLForEnv(envOverride)
	} else {
		wsURL = config.ResolveURL()
	}

	client, err := terminalwire.Connect(wsURL, "micepad", version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Server: %s\n", wsURL)
		if envOverride != "" {
			fmt.Fprintf(os.Stderr, "Environment: %s\n", envOverride)
		}
		fmt.Fprintf(os.Stderr, "Run 'micepad env' to see available environments.\n")
		os.Exit(1)
	}

	if err := client.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// extractEnvFlag pulls -e NAME or --env NAME from the argument list
// and returns the env name plus the remaining args.
func extractEnvFlag(args []string) (string, []string) {
	var envName string
	remaining := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// --env=NAME or -e=NAME
		if strings.HasPrefix(arg, "--env=") {
			envName = strings.TrimPrefix(arg, "--env=")
			continue
		}
		if strings.HasPrefix(arg, "-e=") {
			envName = strings.TrimPrefix(arg, "-e=")
			continue
		}

		// --env NAME or -e NAME
		if (arg == "--env" || arg == "-e") && i+1 < len(args) {
			envName = args[i+1]
			i++ // skip next arg
			continue
		}

		remaining = append(remaining, arg)
	}

	return envName, remaining
}

func handleEnv(args []string) {
	if len(args) == 0 {
		handleEnvList()
		return
	}

	switch args[0] {
	case "list", "ls":
		handleEnvList()
	case "add":
		handleEnvAdd(args[1:])
	case "remove", "rm":
		handleEnvRemove(args[1:])
	case "use":
		handleEnvUse(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown env command: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  micepad env              List environments")
		fmt.Fprintln(os.Stderr, "  micepad env add NAME URL Add an environment")
		fmt.Fprintln(os.Stderr, "  micepad env remove NAME  Remove an environment")
		fmt.Fprintln(os.Stderr, "  micepad env use NAME     Switch active environment")
		os.Exit(1)
	}
}

func handleEnvList() {
	cfg := config.Load()
	fmt.Print(cfg.FormatEnvList())
}

func handleEnvAdd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: micepad env add NAME URL")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, "  micepad env add staging wss://staging.example.com/terminal")
		os.Exit(1)
	}

	name, url := args[0], args[1]
	cfg := config.Load()

	if err := cfg.AddEnv(name, url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added environment %q → %s\n", name, url)
}

func handleEnvRemove(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: micepad env remove NAME")
		os.Exit(1)
	}

	name := args[0]
	cfg := config.Load()

	if err := cfg.RemoveEnv(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed environment %q\n", name)
}

func handleEnvUse(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: micepad env use NAME")
		os.Exit(1)
	}

	name := args[0]
	cfg := config.Load()

	if err := cfg.UseEnv(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	env := cfg.Environments[name]
	fmt.Printf("Switched to %q → %s\n", name, env.URL)
}

func handleConfigure(args []string) {
	cfg := config.Load()
	envName := cfg.CurrentEnv

	// Handle --url flag for non-interactive URL setting (updates current env).
	for i, arg := range args {
		if arg == "--url" && i+1 < len(args) {
			url := args[i+1]
			cfg.Environments[envName] = config.Environment{URL: url}
			if err := cfg.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("URL for %q set to %s\n", envName, url)
			return
		}
	}

	// Interactive configuration — update the current env's URL.
	currentURL := config.ResolveURL()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Server URL for %q [%s]: ", envName, currentURL)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		cfg.Environments[envName] = config.Environment{URL: input}
	}

	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration saved to %s\n", config.Path())
}

func handleVersion() {
	cfg := config.Load()
	fmt.Printf("micepad %s (%s)\n", version, commit)
	fmt.Printf("Env:     %s\n", cfg.CurrentEnv)
	fmt.Printf("Server:  %s\n", config.ResolveURL())
	fmt.Printf("Config:  %s\n", config.Path())
	fmt.Printf("Storage: %s\n", config.Dir())

	if version != "dev" {
		fmt.Println("\nChecking for updates...")
		if latest, err := getLatestVersion(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not check for updates: %v\n", err)
		} else if isNewer(latest, version) {
			fmt.Printf("Update available: v%s → v%s\n", version, latest)
			fmt.Println("Run 'micepad update' to upgrade.")
		} else {
			fmt.Println("Already up to date.")
		}
	}
}

const (
	repo          = "micepad/micepad-cli"
	installScript = "https://github.com/" + repo + "/releases/latest/download/install.sh"
)

// isNewer returns true if latest is a higher semver than current.
func isNewer(latest, current string) bool {
	parse := func(v string) [3]int {
		parts := strings.SplitN(strings.TrimPrefix(v, "v"), ".", 3)
		var nums [3]int
		for i := 0; i < 3 && i < len(parts); i++ {
			nums[i], _ = strconv.Atoi(strings.SplitN(parts[i], "-", 2)[0])
		}
		return nums
	}
	l, c := parse(latest), parse(current)
	for i := 0; i < 3; i++ {
		if l[i] != c[i] {
			return l[i] > c[i]
		}
	}
	return false
}

func getLatestVersion() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get("https://github.com/" + repo + "/releases/latest")
	if err != nil {
		return "", fmt.Errorf("network error (check connectivity)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		return "", fmt.Errorf("unexpected response from GitHub (status %d)", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no redirect from GitHub releases")
	}

	re := regexp.MustCompile(`/v?(\d+\.\d+\.\d+)$`)
	matches := re.FindStringSubmatch(location)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not parse version from redirect URL")
	}
	return matches[1], nil
}

func handleUpdate() {
	fmt.Printf("Current version: %s (%s)\n", version, commit)

	if version != "dev" {
		fmt.Println("Checking for updates...")
		latest, err := getLatestVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not check latest version: %v\n", err)
			fmt.Fprintln(os.Stderr, "Proceeding with update anyway...")
		} else if !isNewer(latest, version) {
			fmt.Printf("Already up to date (v%s)\n", version)
			updateSkill()
			return
		} else {
			fmt.Printf("New version available: v%s → v%s\n", version, latest)
		}
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
