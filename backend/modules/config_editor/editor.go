package config_editor

import (
	"fmt"
	"log"
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

	log.Printf(configPath)

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
