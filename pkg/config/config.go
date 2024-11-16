package config

import (
	"encoding/json"
	"os"

	"github.com/ItzDabbzz/FiveMCarsMerger/pkg/flags"
)

func LoadConfig() (*flags.Flags, error) {
	configPath := "config.json"
	flags := &flags.Flags{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, flags)
	return flags, err
}

func SaveConfig(flags *flags.Flags) error {
	data, err := json.MarshalIndent(flags, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", data, 0644)
}
