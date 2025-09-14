package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lce/backend/modules/app_settings" // Убедитесь, что путь правильный
)

// I18N - это структура, которая будет привязана к фронтенду Wails
type I18N struct {
	// Здесь можно добавить поля, если они нужны для состояния I18N
}

// NewI18N создает новый экземпляр I18N
func NewI18N() *I18N {
	return &I18N{}
}

// GetLanguages возвращает список доступных языков.
// Он сканирует директорию 'locales', читает JSON-файлы
// и извлекает код языка (имя файла) и имя языка (из поля "lang_name" в JSON).
func (i *I18N) GetLanguages() ([]map[string]string, error) { // Изменено на метод
	var langs []map[string]string
	localesDir := "locales" // Предполагаем, что директория 'locales' находится в корне проекта

	entries, err := os.ReadDir(localesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("директория 'locales' не найдена: %w", err)
		}
		return nil, fmt.Errorf("не удалось прочитать директорию 'locales': %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) != ".json" {
			continue
		}

		code := strings.TrimSuffix(filename, filepath.Ext(filename))
		path := filepath.Join(localesDir, filename)

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Ошибка чтения файла '%s': %v\n", path, err)
			continue
		}

		var jsonMap map[string]interface{}
		if err := json.Unmarshal(data, &jsonMap); err != nil {
			fmt.Printf("Ошибка десериализации JSON из файла '%s': %v\n", path, err)
			continue
		}

		name := code // По умолчанию имя равно коду
		if langName, ok := jsonMap["lang_name"].(string); ok {
			name = langName
		}

		langs = append(langs, map[string]string{
			"code": code,
			"name": name,
		})
	}

	return langs, nil
}

// GetCurrentLanguage возвращает текущий выбранный язык из настроек приложения.
func (i *I18N) GetCurrentLanguage() (string, error) { // Изменено на метод
	settings, err := app_settings.LoadSettings()
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить настройки для получения текущего языка: %w", err)
	}
	return settings.Language, nil
}

// SwitchLanguage изменяет текущий язык в настройках приложения и сохраняет их.
func (i *I18N) SwitchLanguage(newLang string) error { // Изменено на метод
	fmt.Printf("Переключение языка на: %s\n", newLang)
	settings, err := app_settings.LoadSettings()
	if err != nil {
		return fmt.Errorf("не удалось загрузить настройки для переключения языка: %w", err)
	}

	settings.Language = newLang

	if err := app_settings.SaveSettings(&settings); err != nil {
		return fmt.Errorf("не удалось сохранить настройки после переключения языка: %w", err)
	}
	return nil
}

// GetTranslationsCurrent возвращает переводы для текущего выбранного языка.
func (i *I18N) GetTranslationsCurrent() (map[string]string, error) { // Изменено на метод
	settings, err := app_settings.LoadSettings()
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить настройки для получения текущих переводов: %w", err)
	}
	langCode := settings.Language
	return i.GetTranslations(langCode) // Вызов метода
}

// GetTranslations возвращает переводы для указанного кода языка.
func (i *I18N) GetTranslations(langCode string) (map[string]string, error) { // Изменено на метод
	path := fmt.Sprintf("locales/%s.json", langCode)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файл не найден, возвращаем пустую карту и ошибку
			return make(map[string]string), fmt.Errorf("файл перевода для '%s' не найден: %w", langCode, err)
		}
		return nil, fmt.Errorf("не удалось прочитать файл перевода '%s': %w", path, err)
	}

	var jsonValue interface{}
	if err := json.Unmarshal(data, &jsonValue); err != nil {
		return nil, fmt.Errorf("не удалось десериализовать JSON перевода из файла '%s': %w", path, err)
	}

	return flattenJSON(jsonValue), nil
}

// flattenJSON рекурсивно преобразует вложенный JSON-объект в плоскую карту.
// Эта функция остается обычной функцией, так как она является внутренней вспомогательной.
func flattenJSON(value interface{}) map[string]string {
	result := make(map[string]string)
	if obj, ok := value.(map[string]interface{}); ok {
		for k, v := range obj {
			key := strings.ToLower(k) // Преобразуем ключ в нижний регистр
			switch val := v.(type) {
			case map[string]interface{}:
				// Рекурсивно обрабатываем вложенные объекты
				nested := flattenJSON(val)
				// Добавляем вложенные ключи без префикса
				for nk, nv := range nested {
					result[nk] = nv
				}
			case string:
				result[key] = val
			default:
				// Игнорируем другие типы значений
			}
		}
	}
	return result
}
