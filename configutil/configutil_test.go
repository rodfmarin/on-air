package configutil

import (
	"os"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	file := "test_config_valid.json"
	content := `{
		"credentials": "credentials.json",
		"token": "token.json",
		"calendar": "primary",
		"days": 7,
		"lifx_token": "token",
		"lifx_light_id": "id",
		"lifx_light_label": "label",
		"lifx_busy_color": "red saturation:0.8",
		"lifx_free_color": "kelvin:3500",
		"reload_interval_seconds": 120
	}`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	defer os.Remove(file)

	cfg, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.CredsPath != "credentials.json" {
		t.Errorf("CredsPath: got %v, want credentials.json", cfg.CredsPath)
	}
	if cfg.TokenPath != "token.json" {
		t.Errorf("TokenPath: got %v, want token.json", cfg.TokenPath)
	}
	if cfg.CalID != "primary" {
		t.Errorf("CalID: got %v, want primary", cfg.CalID)
	}
	if cfg.Days != 7 {
		t.Errorf("Days: got %v, want 7", cfg.Days)
	}
	if cfg.LifxToken != "token" {
		t.Errorf("LifxToken: got %v, want token", cfg.LifxToken)
	}
	if cfg.LifxLightID != "id" {
		t.Errorf("LifxLightID: got %v, want id", cfg.LifxLightID)
	}
	if cfg.LifxLightLabel != "label" {
		t.Errorf("LifxLightLabel: got %v, want label", cfg.LifxLightLabel)
	}
	if cfg.LifxBusyColor != "red saturation:0.8" {
		t.Errorf("LifxBusyColor: got %v, want red saturation:0.8", cfg.LifxBusyColor)
	}
	if cfg.LifxFreeColor != "kelvin:3500" {
		t.Errorf("LifxFreeColor: got %v, want kelvin:3500", cfg.LifxFreeColor)
	}
	if cfg.ReloadIntervalSeconds != 120 {
		t.Errorf("ReloadIntervalSeconds: got %v, want 120", cfg.ReloadIntervalSeconds)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("nonexistent_config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	file := "test_config_invalid.json"
	content := `{"credentials": "credentials.json",` // invalid JSON
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	defer os.Remove(file)

	_, err := LoadConfig(file)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	file := "test_config_empty.json"
	if err := os.WriteFile(file, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	defer os.Remove(file)

	_, err := LoadConfig(file)
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
}

func TestLoadConfig_MissingFields(t *testing.T) {
	file := "test_config_missing_fields.json"
	content := `{"credentials": "credentials.json"}`
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	defer os.Remove(file)

	cfg, err := LoadConfig(file)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.TokenPath != "" {
		t.Errorf("TokenPath: got %v, want empty string", cfg.TokenPath)
	}
	if cfg.Days != 0 {
		t.Errorf("Days: got %v, want 0", cfg.Days)
	}
}
