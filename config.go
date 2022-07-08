package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Directory string `json:"directory"`
}

func LoadOrCreateConfig() (*Config, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("couldn't get configuration: %w", err)
	}

	// We already have the config! Send it back now.
	if config != nil {
		return config, nil
	}

	var directory string
	if err := survey.AskOne(&survey.Input{
		Message: "Directory to store reproductions",
		Suggest: func(toComplete string) []string {
			if strings.Contains(toComplete, "~") {
				toComplete, _ = homedir.Expand(toComplete)
			}

			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &directory); err != nil {
		return nil, fmt.Errorf("couldn't get answers: %w", err)
	}

	path, err := homedir.Expand(directory)
	if err != nil {
		return nil, fmt.Errorf("couldn't expand directory filepath: %w", err)
	}

	// check if the directory exists.
	if _, err := os.Stat(path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("couldn't stat directory: %w", err)
		}

		var create bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Directory doesn't exist, create it?",
			Default: true,
		}, &create); err != nil {
			return nil, fmt.Errorf("couldn't get answers: %w", err)
		}

		if create {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return nil, fmt.Errorf("couldn't create directory %s: %w", path, err)
			}
		}
	}

	config = &Config{
		Directory: path,
	}

	if err := SaveConfig(config); err != nil {
		return nil, fmt.Errorf("couldnt' save configuration: %w", err)
	}

	return config, nil
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't get home directory: %w", err)
	}

	path := filepath.Join(home, ".reprc.json")

	return path, nil
}

func LoadConfig() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get config path: %w", err)
	}

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("could not open configuration at %s: %w", path, err)
	}

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode configuration at %s: %w", path, err)
	}

	// Always expand the directory, it could have been user edited.
	directory, err := homedir.Expand(config.Directory)
	if err != nil {
		return nil, fmt.Errorf("couldn't expand directory filepath: %w", err)
	}
	config.Directory = directory

	return &config, nil
}

func SaveConfig(config *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return fmt.Errorf("couldn't get config path: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create configuration at %s: %w", path, err)
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("could not encode configuration at %s: %w", path, err)
	}

	fmt.Printf("Configuration saved %s\n", path)

	return nil
}
