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

// Read loads the config file from the home directory and returns it.
func Read() (Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(homedir, ConfigFileName)

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil

}

// SetUser updates CurrentUserName in memory and rewrites the config file.
// os.Create truncates the file, so the whole struct is written fresh each time.
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homedir, ConfigFileName)

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
