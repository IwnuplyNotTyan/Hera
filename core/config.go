package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds persisted user settings (theme, centering, background fill, language).
// JSON tags map to fields in ~/.config/hera/config.json.
type Config struct {
	ThemeName    string `json:"theme"`
	CenterWindow bool   `json:"centered"`
	Background   bool   `json:"background"`
	Language     string `json:"language,omitempty"`
}

// InitConfig loads settings from ~/.config/hera/config.json.
// If the file does not exist it is created, using m.Config as defaults
// (or hardcoded defaults when m.Config is nil). Returns an error when the
// config directory cannot be determined or created, or when the file cannot
// be written.
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

	cfg := m.Config
	if cfg == nil {
		lang := "en"
		if m != nil {
			lang = m.Localizer.GetLanguage()
		}
		cfg = &Config{
			ThemeName:    "dark",
			CenterWindow: true,
			Background:   false,
			Language:     lang,
		}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	m.Config = cfg
	return nil
}

// SaveConfig writes the current Config to ~/.config/hera/config.json.
// Returns nil when m.Config is nil. Returns an error when the config
// directory cannot be determined, marshalling fails, or the file cannot
// be written.
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

// SetLanguage switches the UI language at runtime, persists the choice to
// config, and rebuilds key bindings so help text matches the new language.
func (m *Model) SetLanguage(lang string) error {
	if err := m.Localizer.SetLanguage(lang); err != nil {
		return err
	}
	m.keys = newKeyMap(m.Localizer)
	m.menuKeys = newMenuKeyMap(m.Localizer)
	if m.Config != nil {
		m.Config.Language = lang
		if err := m.SaveConfig(); err != nil {
			return err
		}
	}
	return nil
}
