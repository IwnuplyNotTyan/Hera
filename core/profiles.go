package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (m *Model) profilesPath() (string, error) {
	d, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "profiles.json"), nil
}

// SaveProfiles writes the profiles array to the profiles.json config file.
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

// LoadProfiles reads the profiles array from the profiles.json config file.
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
