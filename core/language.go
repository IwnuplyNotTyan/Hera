package generate

// SetLanguage changes the active locale and updates keybindings and config.
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
			if setErr := m.Localizer.SetLanguage(prevLang); setErr != nil {
				return setErr
			}
			m.keys = newKeyMap(m.Localizer)
			m.menuKeys = newMenuKeyMap(m.Localizer)
			m.Config.Language = prevLang
			return err
		}
	}
	return nil
}
