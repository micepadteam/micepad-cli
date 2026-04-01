package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const DefaultURL = "wss://studio.micepad.co/terminal"

// Built-in environments available out of the box.
var BuiltinEnvironments = map[string]Environment{
	"prod":  {URL: "wss://studio.micepad.co/terminal"},
	"alpha": {URL: "wss://launchpad.micepad.co/terminal"},
	"dev":   {URL: "ws://localhost:3000/terminal"},
}

const DefaultEnv = "prod"

// Environment holds the configuration for a single named environment.
type Environment struct {
	URL string `json:"url"`
}

// Config holds persistent CLI configuration stored at ~/.micepad/config.json.
type Config struct {
	CurrentEnv   string                 `json:"current_env,omitempty"`
	Environments map[string]Environment `json:"environments,omitempty"`

	// Legacy field — migrated to Environments on load.
	LegacyURL string `json:"url,omitempty"`
}

// Dir returns the micepad config directory (~/.micepad).
func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".micepad")
}

// Path returns the config file path (~/.micepad/config.json).
func Path() string {
	return filepath.Join(Dir(), "config.json")
}

// Load reads the config file and migrates legacy format if needed.
func Load() *Config {
	cfg := &Config{}
	data, err := os.ReadFile(Path())
	if err != nil {
		// No config file — seed with built-in environments.
		cfg.seed()
		return cfg
	}
	json.Unmarshal(data, cfg)
	cfg.migrate()
	return cfg
}

// seed populates a fresh config with built-in environments.
func (c *Config) seed() {
	c.CurrentEnv = DefaultEnv
	c.Environments = make(map[string]Environment)
	for name, env := range BuiltinEnvironments {
		c.Environments[name] = env
	}
}

// migrate converts legacy single-URL config to the environments format.
func (c *Config) migrate() {
	if c.Environments == nil {
		c.Environments = make(map[string]Environment)
	}

	// Seed any missing built-in environments.
	for name, env := range BuiltinEnvironments {
		if _, exists := c.Environments[name]; !exists {
			c.Environments[name] = env
		}
	}

	// Migrate legacy "url" field.
	if c.LegacyURL != "" && c.CurrentEnv == "" {
		// Find if the legacy URL matches a built-in env.
		matched := false
		for name, env := range BuiltinEnvironments {
			if env.URL == c.LegacyURL {
				c.CurrentEnv = name
				matched = true
				break
			}
		}
		if !matched {
			// Save as a custom environment.
			c.Environments["custom"] = Environment{URL: c.LegacyURL}
			c.CurrentEnv = "custom"
		}
		c.LegacyURL = ""
	}

	if c.CurrentEnv == "" {
		c.CurrentEnv = DefaultEnv
	}
}

// Save writes the config to disk, creating the directory if needed.
func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return err
	}
	// Clear legacy field on save.
	c.LegacyURL = ""
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0644)
}

// ResolveURL returns the server URL using priority:
// MICEPAD_URL env > current environment > default.
func ResolveURL() string {
	if envURL := os.Getenv("MICEPAD_URL"); envURL != "" {
		return envURL
	}
	return ResolveURLForEnv("")
}

// ResolveURLForEnv returns the URL for a specific environment name.
// If envName is empty, uses the current environment.
func ResolveURLForEnv(envName string) string {
	if envURL := os.Getenv("MICEPAD_URL"); envURL != "" && envName == "" {
		return envURL
	}
	cfg := Load()
	if envName == "" {
		envName = cfg.CurrentEnv
	}
	if env, ok := cfg.Environments[envName]; ok {
		return env.URL
	}
	return DefaultURL
}

// AddEnv adds a named environment.
func (c *Config) AddEnv(name, url string) error {
	if name == "" {
		return fmt.Errorf("environment name cannot be empty")
	}
	c.Environments[name] = Environment{URL: url}
	return nil
}

// RemoveEnv removes a named environment. Cannot remove the active environment.
func (c *Config) RemoveEnv(name string) error {
	if name == c.CurrentEnv {
		return fmt.Errorf("cannot remove active environment %q — switch first with: micepad env use <name>", name)
	}
	if _, ok := c.Environments[name]; !ok {
		return fmt.Errorf("environment %q not found", name)
	}
	delete(c.Environments, name)
	return nil
}

// UseEnv switches the active environment.
func (c *Config) UseEnv(name string) error {
	if _, ok := c.Environments[name]; !ok {
		return fmt.Errorf("environment %q not found — add it with: micepad env add %s <url>", name, name)
	}
	c.CurrentEnv = name
	return nil
}

// EnvNames returns sorted environment names.
func (c *Config) EnvNames() []string {
	names := make([]string, 0, len(c.Environments))
	for name := range c.Environments {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// FormatEnvList returns a formatted string listing all environments.
func (c *Config) FormatEnvList() string {
	var b strings.Builder
	for _, name := range c.EnvNames() {
		env := c.Environments[name]
		marker := "  "
		if name == c.CurrentEnv {
			marker = "★ "
		}
		fmt.Fprintf(&b, "%s%-8s %s\n", marker, name, env.URL)
	}
	return b.String()
}
