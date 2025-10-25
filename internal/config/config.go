package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigPath()

	if err != nil {
		return Config{}, err
	}

	fileData, err := os.ReadFile(configPath)

	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(fileData, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}


func getConfigPath() (string, error) {
	// func returns a string with the path to the config file in the user's home directory
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	// will return /[your-home-directory]/.gatorconfig.json
	return filepath.Join(homeDir, configFilename), nil
}

func write(config Config) error {
	configPath, err := getConfigPath()

	if err != nil {
		return err
	}

	fileData, err := json.Marshal(config)

	if err != nil {
		return err
	}

	return os.WriteFile(configPath, fileData, 0644)
}