package app_settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync" // Для обеспечения потокобезопасности при изменении настроек
)

// Settings - структура для хранения настроек приложения
type Settings struct {
	Width    uint32   `json:"width"`
	Height   uint32   `json:"height"`
	Language string   `json:"language"`
	GamePath string   `json:"game_path"`
	FirstRun bool     `json:"first_run"` // NEW: Добавляем поле FirstRun
	AllPaths []string `json:"all_paths"` // NEW: Добавляем поле AllPaths
	Theme    string   `json:"theme"`     // NEW: Добавляем поле Theme
}

// DefaultSettings возвращает настройки по умолчанию
func DefaultSettings() Settings {
	return Settings{
		Width:    1600,
		Height:   900,
		Language: "en",
		GamePath: "",
		FirstRun: true,       // NEW: Значение по умолчанию для FirstRun
		AllPaths: []string{}, // NEW: Значение по умолчанию для AllPaths
		Theme:    "default",  // NEW: Значение по умолчанию для Theme
	}
}

// getSettingsPath возвращает путь к файлу настроек
func getSettingsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("не удалось получить директорию конфигурации пользователя: %w", err)
	}

	appConfigDir := filepath.Join(configDir, "LoD Config Editor")
	if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
		if err := os.MkdirAll(appConfigDir, 0755); err != nil {
			return "", fmt.Errorf("не удалось создать директорию конфигурации приложения: %w", err)
		}
	}

	return filepath.Join(appConfigDir, "settings.json"), nil
}

// LoadSettings загружает настройки из файла
func LoadSettings() (Settings, error) {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return DefaultSettings(), err
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл не существует, сохраняем настройки по умолчанию
			defaultSettings := DefaultSettings()
			if err := SaveSettings(&defaultSettings); err != nil {
				return defaultSettings, fmt.Errorf("не удалось сохранить настройки по умолчанию: %w", err)
			}
			return defaultSettings, nil
		}
		return DefaultSettings(), fmt.Errorf("не удалось прочитать файл настроек: %w", err)
	}

	var userSettings Settings
	// Инициализируем userSettings значениями по умолчанию, чтобы отсутствующие поля сохраняли их
	defaultValues := DefaultSettings()
	userSettings = defaultValues

	if err := json.Unmarshal(data, &userSettings); err != nil {
		// Ошибка десериализации, используем настройки по умолчанию и логируем ошибку
		fmt.Printf("Ошибка десериализации настроек, используются настройки по умолчанию: %v\n", err)
		// Пытаемся сохранить настройки по умолчанию, если десериализация не удалась, чтобы файл был валидным
		if err := SaveSettings(&defaultValues); err != nil {
			fmt.Printf("Не удалось сохранить настройки по умолчанию после ошибки десериализации: %v\n", err)
		}
		return defaultValues, fmt.Errorf("ошибка десериализации: %w", err)
	}

	// Минимальные размеры
	if userSettings.Width < 650 {
		userSettings.Width = 650
	}
	if userSettings.Height < 350 {
		userSettings.Height = 350
	}

	// Проверяем и добавляем отсутствующие ключи (если структура изменилась или были пустые значения)
	updated := false
	if userSettings.Language == "" {
		userSettings.Language = defaultValues.Language
		updated = true
	}
	if userSettings.GamePath == "" {
		userSettings.GamePath = defaultValues.GamePath
		updated = true
	}
	if userSettings.AllPaths == nil { // Проверяем на nil, если поле отсутствовало в старых настройках
		userSettings.AllPaths = defaultValues.AllPaths
		updated = true
	}
	if userSettings.Theme == "" { // Проверяем на пустую строку, если поле отсутствовало
		userSettings.Theme = defaultValues.Theme
		updated = true
	}

	if updated {
		if err := SaveSettings(&userSettings); err != nil {
			return userSettings, fmt.Errorf("не удалось сохранить обновленные настройки: %w", err)
		}
	}

	return userSettings, nil
}

// SaveSettings сохраняет настройки в файл
func SaveSettings(newSettings *Settings) error {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(newSettings, "", "    ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать настройки в JSON: %w", err)
	}

	if err := os.WriteFile(settingsPath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("не удалось записать файл настроек: %w", err)
	}
	return nil
}

// AppSettings - это структура, которая будет привязана к фронтенду Wails.
// Она содержит текущие настройки и мьютекс для потокобезопасного доступа.
type AppSettings struct {
	settings Settings
	lock     sync.RWMutex // Мьютекс для безопасного доступа к настройкам
}

// NewAppSettings создает новый экземпляр AppSettings.
// Он загружает существующие настройки или использует значения по умолчанию.
func NewAppSettings() *AppSettings {
	s, err := LoadSettings()
	if err != nil {
		fmt.Printf("Ошибка при загрузке настроек: %v. Используются настройки по умолчанию.\n", err)
		s = DefaultSettings()
	}
	return &AppSettings{
		settings: s,
	}
}

// GetSettings возвращает текущие настройки приложения.
// Эта функция привязана к фронтенду Wails.
func (a *AppSettings) GetSettings() Settings {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.settings
}

// UpdateSettings обновляет настройки приложения на основе предоставленной карты.
// Эта функция привязана к фронтенду Wails.
func (a *AppSettings) UpdateSettings(newSettings map[string]interface{}) (Settings, error) {
	fmt.Println("=== Начало UpdateSettings ===")
	fmt.Printf("Входные данные: %+v\n", newSettings)

	a.lock.Lock()
	defer a.lock.Unlock()

	currentSettings := a.settings // Создаем копию для работы

	fmt.Printf("Текущие настройки перед обновлением: %+v\n", currentSettings)

	updated := false
	for key, value := range newSettings {
		fmt.Printf("Обработка поля '%s': %+v\n", key, value)
		switch key {
		case "width":
			if v, ok := value.(float64); ok { // Числа из JSON - это float64 в Go
				currentSettings.Width = uint32(v)
				updated = true
			}
		case "height":
			if v, ok := value.(float64); ok {
				currentSettings.Height = uint32(v)
				updated = true
			}
		case "language":
			if v, ok := value.(string); ok {
				currentSettings.Language = v
				updated = true
			}
		case "game_path":
			if v, ok := value.(string); ok {
				currentSettings.GamePath = v
				updated = true
			}
		case "first_run": // NEW: Добавляем обработку для first_run
			if v, ok := value.(bool); ok {
				currentSettings.FirstRun = v
				updated = true
			}
		case "all_paths": // NEW: Добавляем обработку для all_paths
			if v, ok := value.([]interface{}); ok { // JSON массив будет []interface{}
				var paths []string
				for _, item := range v {
					if pathStr, ok := item.(string); ok {
						paths = append(paths, pathStr)
					}
				}
				currentSettings.AllPaths = paths
				updated = true
			}
		case "theme": // NEW: Добавляем обработку для theme
			if v, ok := value.(string); ok {
				currentSettings.Theme = v
				updated = true
			}
		default:
			fmt.Printf("Неизвестное поле: %s\n", key)
		}
	}

	// Повторная проверка минимальных размеров после обновлений
	if currentSettings.Width < 650 {
		currentSettings.Width = 650
		updated = true
	}
	if currentSettings.Height < 350 {
		currentSettings.Height = 350
		updated = true
	}

	if updated {
		if err := SaveSettings(&currentSettings); err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			return a.settings, fmt.Errorf("не удалось сохранить настройки: %w", err) // Возвращаем старые настройки в случае ошибки
		}
		a.settings = currentSettings // Обновляем внутренние настройки только если сохранение было успешным
		fmt.Println("Настройки успешно сохранены")
	} else {
		fmt.Println("Настройки не изменились, сохранение не требуется.")
	}

	fmt.Println("=== Конец UpdateSettings ===")
	return a.settings, nil
}

// GetOption возвращает значение определенной опции по ключу.
// Эта функция привязана к фронтенду Wails.
func (a *AppSettings) GetOption(key string) interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()

	switch key {
	case "width":
		return a.settings.Width
	case "height":
		return a.settings.Height
	case "language":
		return a.settings.Language
	case "game_path":
		return a.settings.GamePath
	case "first_run": // NEW: Добавляем возврат значения для first_run
		return a.settings.FirstRun
	case "all_paths": // NEW: Добавляем возврат значения для all_paths
		return a.settings.AllPaths
	case "theme": // NEW: Добавляем возврат значения для theme
		return a.settings.Theme
	default:
		return nil
	}
}
