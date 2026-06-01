package i18n

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Localizer interface {
	T(key string, args ...interface{}) string
	GetLanguage() string
	SetLanguage(lang string) error
	AvailableLanguages() []string
	RandomEasterEgg() string
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

	if err := t.loadEmbedded(); err == nil {
		if err := t.validateLoadedLocales(); err != nil {
			return nil, err
		}
		return t, nil
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

	if err := t.validateLoadedLocales(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Translator) validateLoadedLocales() error {
	if len(t.translations) == 0 {
		return fmt.Errorf("no locales loaded")
	}
	if _, ok := t.translations[t.defaultLang]; !ok {
		return fmt.Errorf("default language %q not found in loaded locales", t.defaultLang)
	}
	return nil
}

func (t *Translator) loadEmbedded() error {
	entries, err := embeddedLocales.ReadDir("locales")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		lang := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		data, err := embeddedLocales.ReadFile(path.Join("locales", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read embedded locale %s: %w", entry.Name(), err)
		}

		var translations map[string]interface{}
		if err := json.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse embedded locale %s: %w", entry.Name(), err)
		}

		t.translations[lang] = flattenMap(translations)
	}

	if len(t.translations) == 0 {
		return fmt.Errorf("no embedded locales found")
	}

	return nil
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
	if len(args) == 0 {
		return template
	}

	result := template

	switch firstArg := args[0].(type) {
	case map[string]interface{}:
		for key, val := range firstArg {
			result = strings.ReplaceAll(result, "{"+key+"}", fmt.Sprintf("%v", val))
		}
		for i := 1; i < len(args); i++ {
			placeholder := fmt.Sprintf("{%d}", i)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", args[i]))
		}
	case int:
		result = strings.ReplaceAll(result, "{n}", fmt.Sprintf("%d", firstArg))
		result = strings.ReplaceAll(result, "{0}", fmt.Sprintf("%d", firstArg))
		if len(args) > 1 {
			result = strings.ReplaceAll(result, "{hp}", fmt.Sprintf("%v", args[1]))
		}
		for i := 1; i < len(args); i++ {
			placeholder := fmt.Sprintf("{%d}", i)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", args[i]))
		}
	default:
		result = strings.ReplaceAll(result, "{n}", fmt.Sprintf("%v", firstArg))
		result = strings.ReplaceAll(result, "{0}", fmt.Sprintf("%v", firstArg))
		if len(args) > 1 {
			result = strings.ReplaceAll(result, "{hp}", fmt.Sprintf("%v", args[1]))
		}
		for i := 1; i < len(args); i++ {
			placeholder := fmt.Sprintf("{%d}", i)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", args[i]))
		}
	}

	return result
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

func (t *Translator) RandomEasterEgg() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	eggs, ok := t.translations[t.currentLang]["easterEggs"].([]any)
	if !ok || len(eggs) == 0 {
		eggs, ok = t.translations[t.defaultLang]["easterEggs"].([]any)
		if !ok || len(eggs) == 0 {
			return "???（ ╯°□°）╯︵ ┻━┻"
		}
	}

	idx := rand.Intn(len(eggs))
	egg, ok := eggs[idx].(string)
	if !ok {
		return "???（ ╯°□°）╯︵ ┻━┻"
	}
	return egg
}
