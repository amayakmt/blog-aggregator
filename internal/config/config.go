package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

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
