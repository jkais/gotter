package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	Timestamp int64 `json:"timestamp"`
	Response  map[string]any `json:"response"`
}

func getCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gotter", "cache.json")
}

func ensureCacheFile() {
	cachePath := getCachePath()
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		cache := Cache{Timestamp: 0, Response: make(map[string]any)}
		data, _ := json.MarshalIndent(cache, "", "  ")
		os.WriteFile(cachePath, data, 0644)
	}
}

func getCacheData() Cache {
	ensureCacheFile()
	cachePath := getCachePath()
	file, _ := os.Open(cachePath)
	defer file.Close()
	var cache Cache
	json.NewDecoder(file).Decode(&cache)
	return cache
}

func needsUpdate() bool {
	return time.Now().Unix() - getCacheData().Timestamp > 600
}

func setCache(jsonResponse string) {
	var responseData map[string]any
	json.Unmarshal([]byte(jsonResponse), &responseData)
	prettyResponse, _ := json.MarshalIndent(responseData, "  ", "  ")
	cache := fmt.Sprintf(`{
		"timestamp": %d,
		"response": %s
		}`, time.Now().Unix(), string(prettyResponse))
	os.WriteFile(getCachePath(), []byte(cache), 0644)
}

func getCache() map[string]any {
	cache := getCacheData()
	return cache.Response
}
