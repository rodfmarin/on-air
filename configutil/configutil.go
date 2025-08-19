// Package configutil provides configuration loading for the app.
package configutil

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	CredsPath             string `json:"credentials"`
	TokenPath             string `json:"token"`
	CalID                 string `json:"calendar"`
	Days                  int    `json:"days"`
	LifxToken             string `json:"lifx_token"`
	LifxLightID           string `json:"lifx_light_id"`
	LifxLightLabel        string `json:"lifx_light_label"`
	LifxBusyColor         string `json:"lifx_busy_color"`
	LifxFreeColor         string `json:"lifx_free_color"`
	ReloadIntervalSeconds int    `json:"reload_interval_seconds"`
}

// LoadConfig loads config from the given file path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {

		}
	}(f)
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}
