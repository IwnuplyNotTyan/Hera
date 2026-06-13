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
	CenterWindow string `json:"centered"`
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
		if m.Localizer != nil {
			lang = m.Localizer.GetLanguage()
		}
		cfg = &Config{
			ThemeName:    "dark",
			CenterWindow: "c",
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

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Migration: convert old bool "centered" to new string format
	if centered, ok := raw["centered"]; ok {
		switch v := centered.(type) {
		case bool:
			if v {
				raw["centered"] = "c"
			} else {
				raw["centered"] = "tl"
			}
		}
	}

	// Ensure default if missing
	if _, ok := raw["centered"]; !ok {
		raw["centered"] = "c"
	}

	fixed, err := json.Marshal(raw)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(fixed, &cfg); err != nil {
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

func (m *Model) profilesPath() (string, error) {
	d, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "profiles.json"), nil
}

// SaveProfiles persists the profile array to ~/.config/hera/profiles.json.
func (m *Model) SaveProfiles() error {
	p, err := m.profilesPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(m.Profiles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

// LoadProfiles reads the profile array from ~/.config/hera/profiles.json.
func (m *Model) LoadProfiles() error {
	p, err := m.profilesPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.Profiles)
}

// SetLanguage switches the UI language at runtime, persists the choice to
// config, and rebuilds key bindings so help text matches the new language.
// On SaveConfig failure the previous language is restored.
func (m *Model) SetLanguage(lang string) error {
	prevLang := m.Localizer.GetLanguage()
	if err := m.Localizer.SetLanguage(lang); err != nil {
		return err
	}
	m.keys = newKeyMap(m.Localizer)
	m.menuKeys = newMenuKeyMap(m.Localizer)
	if m.Config != nil {
		m.Config.Language = lang
		if err := m.SaveConfig(); err != nil {
			_ = m.Localizer.SetLanguage(prevLang)
			m.keys = newKeyMap(m.Localizer)
			m.menuKeys = newMenuKeyMap(m.Localizer)
			m.Config.Language = prevLang
			return err
		}
	}
	return nil
}
