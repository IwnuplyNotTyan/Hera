package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ThemeName  string `json:"theme"`
	CenterWindow   bool   `json:"centered"`
	Background bool   `json:"background"`
}

func (m *Model) InitConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configPath); err == nil {
		return m.loadConfig(configPath)
	}

	defaultConfig := &Config{
		ThemeName:    "dark",
		CenterWindow: true,
		Background:   false,
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	m.Config = defaultConfig
	return nil
}

func (m *Model) SaveConfig() error {
	if m.Config == nil {
		return nil
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(m.Config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (m *Model) loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	m.Config = &cfg
	return nil
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "hera"), nil
}
