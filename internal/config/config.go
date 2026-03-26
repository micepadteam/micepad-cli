package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const DefaultURL = "wss://studio.micepad.co/terminal"

// Config holds persistent CLI configuration stored at ~/.micepad/config.json.
type Config struct {
	URL string `json:"url,omitempty"`
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

// Load reads the config file. Returns an empty config if the file doesn't exist.
func Load() *Config {
	cfg := &Config{}
	data, err := os.ReadFile(Path())
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, cfg)
	return cfg
}

// Save writes the config to disk, creating the directory if needed.
func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0644)
}

// ResolveURL returns the server URL using priority: MICEPAD_URL env > config file > default.
func ResolveURL() string {
	if envURL := os.Getenv("MICEPAD_URL"); envURL != "" {
		return envURL
	}
	cfg := Load()
	if cfg.URL != "" {
		return cfg.URL
	}
	return DefaultURL
}
