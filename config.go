package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ApiKey string `json:"api_key"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gotter", "config.json")
}

func loadConfig() *Config {
	file, _ := os.Open(getConfigPath())
	defer file.Close()
	var config Config
	json.NewDecoder(file).Decode(&config)
	return &config
}
