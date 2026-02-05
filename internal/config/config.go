package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	TasksDir   = ".tasks"
	ConfigFile = "config.toml"
)

type Config struct {
	Board BoardConfig `toml:"board"`
}

type BoardConfig struct {
	Columns []string `toml:"columns"`
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Board: BoardConfig{
			Columns: []string{"TODO", "DOING", "DONE"},
		},
	}
}

// Load reads and parses the config file from the given path
func Load(path string) (*Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to the given path
func Save(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// Initialize ensures .tasks/ directory and config.toml exist
// Returns true if new config was created, false if it already existed
func Initialize() (bool, error) {
	// Create .tasks directory if it doesn't exist
	if err := os.MkdirAll(TasksDir, 0755); err != nil {
		return false, fmt.Errorf("failed to create %s directory: %w", TasksDir, err)
	}

	configPath := filepath.Join(TasksDir, ConfigFile)

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		// Config exists, don't overwrite
		return false, nil
	}

	// Create default config
	defaultCfg := GetDefaultConfig()
	if err := Save(configPath, defaultCfg); err != nil {
		return false, err
	}

	return true, nil
}
