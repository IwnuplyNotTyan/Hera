package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Localizer interface {
	T(key string, args ...interface{}) string
	GetLanguage() string
	SetLanguage(lang string) error
	AvailableLanguages() []string
}

type Translator struct {
	translations map[string]map[string]interface{}
	currentLang  string
	defaultLang  string
	mu           sync.RWMutex
}

func NewTranslator(localesPath string, defaultLang string) (*Translator, error) {
	t := &Translator{
		translations: make(map[string]map[string]interface{}),
		defaultLang:  defaultLang,
		currentLang:  defaultLang,
	}

	files, err := os.ReadDir(localesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		lang := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		data, err := os.ReadFile(filepath.Join(localesPath, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read locale file %s: %w", file.Name(), err)
		}

		var translations map[string]interface{}
		if err := json.Unmarshal(data, &translations); err != nil {
			return nil, fmt.Errorf("failed to parse locale file %s: %w", file.Name(), err)
		}

		t.translations[lang] = flattenMap(translations)
	}

	return t, nil
}

func flattenMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	flattenRecursive(m, "", result)
	return result
}

func flattenRecursive(m map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]interface{}); ok {
			flattenRecursive(nested, fullKey, result)
		} else {
			result[fullKey] = value
		}
	}
}

func (t *Translator) SetLanguage(lang string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.translations[lang]; !ok {
		return fmt.Errorf("language %s not found", lang)
	}

	t.currentLang = lang
	return nil
}

func (t *Translator) GetLanguage() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.currentLang
}

func (t *Translator) T(key string, args ...interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	translations := t.translations[t.currentLang]
	if translations == nil {
		translations = t.translations[t.defaultLang]
	}

	value, ok := translations[key]
	if !ok {
		value, ok = t.translations[t.defaultLang][key]
	}

	if !ok {
		return key
	}

	template, ok := value.(string)
	if !ok {
		return key
	}

	if len(args) == 0 {
		return template
	}

	return t.interpolate(template, args...)
}

func (t *Translator) interpolate(template string, args ...interface{}) string {
	parts := strings.Split(template, "{")

	if len(parts) == 1 {
		return template
	}

	var result strings.Builder
	result.WriteString(parts[0])

	for i := 1; i < len(parts); i++ {
		closing := strings.Index(parts[i], "}")
		if closing == -1 {
			result.WriteString("{")
			result.WriteString(parts[i])
			continue
		}

		key := parts[i][:closing]
		rest := parts[i][closing+1:]

		found := false
		for j, arg := range args {
			switch v := arg.(type) {
			case map[string]interface{}:
				if val, ok := v[key]; ok {
					result.WriteString(fmt.Sprintf("%v", val))
					found = true
					break
				}
			default:
				if key == "n" || key == "x" || key == "y" {
					result.WriteString(fmt.Sprintf("%v", v))
					found = true
					break
				}
				if j == 0 {
					result.WriteString(fmt.Sprintf("%v", v))
					found = true
					break
				}
			}
		}

		if !found {
			result.WriteString("{")
			result.WriteString(key)
		}

		result.WriteString(rest)
	}

	return result.String()
}

func (t *Translator) AvailableLanguages() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	languages := make([]string, 0, len(t.translations))
	for lang := range t.translations {
		languages = append(languages, lang)
	}
	return languages
}
