package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ConfigFileName is the JSON file stored in the user's home directory that
// persists the active username across CLI invocations.
const ConfigFileName = ".gatorconfig.json"

// Config is the on-disk representation of the application's settings.
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// configFilePath returns the absolute path to the config file.
func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigFileName), nil
}

// Read loads the config file from the home directory and returns it.
func Read() (Config, error) {
	configPath, err := configFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	cfg := Config{}
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// SetUser updates CurrentUserName in memory and rewrites the config file.
// os.Create truncates the file, so the whole struct is written fresh each time.
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	configPath, err := configFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cfg)
}
