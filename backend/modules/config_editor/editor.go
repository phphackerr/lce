package config_editor

import (
	"fmt"
	"strings"
	"sync"

	"lce/backend/modules/app_settings"
)

type ConfigEditor struct {
	mu     sync.Mutex
	config *GameConfig
}

func NewConfigEditor() *ConfigEditor {
	return &ConfigEditor{
		config: &GameConfig{},
	}
}

// Загрузить конфиг
func (e *ConfigEditor) LoadConfig() error {
	settings, err := app_settings.LoadSettings()
	if err != nil {
		return fmt.Errorf("ошибка загрузки настроек: %w", err)
	}

	wcDir := settings.GamePath
	configPath := wcDir + "\\config.lod.ini"

	e.mu.Lock()
	defer e.mu.Unlock()
	return e.config.Load(configPath)
}

// Проверить наличие
func (e *ConfigEditor) IsConfigAvailable() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.config != nil && e.config.Get("", "") != ""
}

// Получить значение
func (e *ConfigEditor) GetConfigValue(section, option, defaultValue string) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	val := e.config.Get(section, option)
	if val == "" {
		return defaultValue
	}
	return val
}

// Установить значение
func (e *ConfigEditor) SetConfigValue(section, option, value string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.config.Set(section, option, value)
	return e.config.Save()
}

// Перезагрузить
func (e *ConfigEditor) ReloadConfig() error {
	e.mu.Lock()
	path := e.config.Path()
	e.mu.Unlock()
	return e.config.Load(path)
}

// Получить значение как Hotkey
func (e *ConfigEditor) GetHotkeyValue(section, option string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	rawValue := e.config.Get(section, option)
	if rawValue == "" {
		return "", nil
	}

	// Если ctrl / shift / alt
	if rawValue == "ctrl" || rawValue == "shift" || rawValue == "alt" {
		return string(rawValue[0]-32) + rawValue[1:], nil // Ctrl / Shift / Alt
	}

	// Если начинается с 0x
	if len(rawValue) > 2 && rawValue[:2] == "0x" {
		keyName := Lookup(rawValue)
		if keyName == "" {
			return rawValue, nil
		}
		// F-клавиши
		if keyName[0] == 'f' && len(keyName) > 1 {
			return "F" + keyName[1:], nil
		}
		// Одиночные буквы
		if len(keyName) == 1 {
			return strings.ToUpper(keyName), nil
		}
		return keyName, nil
	}

	return rawValue, nil
}

// CheckConfigDiff сравнивает текущий конфиг в памяти с тем, что на диске,
// возвращает map[section]map[key]value только с изменёнными значениями
func (e *ConfigEditor) CheckConfigDiff() (map[string]map[string]map[string]string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.config == nil || e.config.Path() == "" {
		return nil, fmt.Errorf("config not loaded")
	}

	// Загружаем версию с диска
	diskCfg := &GameConfig{}
	if err := diskCfg.Load(e.config.Path()); err != nil {
		return nil, fmt.Errorf("failed to load config from disk: %w", err)
	}

	diff := make(map[string]map[string]map[string]string)

	// --- Проверяем новые и изменённые секции/ключи ---
	for _, section := range diskCfg.file.Sections() {
		secName := section.Name()
		memSection := e.config.file.Section(secName)

		if memSection == nil || memSection.KeysHash() == nil {
			// Секция полностью новая
			diff[secName] = make(map[string]map[string]string)
			diff[secName]["<section>"] = map[string]string{"old": "", "new": "<added>", "status": "added"}
			for _, key := range section.Keys() {
				diff[secName][key.Name()] = map[string]string{
					"old":    "",
					"new":    key.Value(),
					"status": "added",
				}
			}
			continue
		}

		// Проверяем ключи
		for _, key := range section.Keys() {
			diskVal := key.Value()
			if memSection.HasKey(key.Name()) {
				memVal := memSection.Key(key.Name()).String()
				if diskVal != memVal {
					if diff[secName] == nil {
						diff[secName] = make(map[string]map[string]string)
					}
					diff[secName][key.Name()] = map[string]string{
						"old":    memVal,
						"new":    diskVal,
						"status": "modified",
					}
				}
			} else {
				if diff[secName] == nil {
					diff[secName] = make(map[string]map[string]string)
				}
				diff[secName][key.Name()] = map[string]string{
					"old":    "",
					"new":    diskVal,
					"status": "added",
				}
			}
		}
	}

	// --- Проверяем удалённые секции/ключи ---
	for _, section := range e.config.file.Sections() {
		secName := section.Name()
		diskSection := diskCfg.file.Section(secName)
		if diskSection == nil || diskSection.KeysHash() == nil {
			// Секция полностью удалена
			diff[secName] = make(map[string]map[string]string)
			diff[secName]["<section>"] = map[string]string{
				"old":    "<deleted>",
				"new":    "",
				"status": "deleted",
			}
			continue
		}

		for _, key := range section.Keys() {
			if !diskSection.HasKey(key.Name()) {
				if diff[secName] == nil {
					diff[secName] = make(map[string]map[string]string)
				}
				diff[secName][key.Name()] = map[string]string{
					"old":    key.Value(),
					"new":    "",
					"status": "deleted",
				}
			}
		}
	}

	return diff, nil
}
