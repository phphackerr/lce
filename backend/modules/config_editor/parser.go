package config_editor

import (
	"log"

	"gopkg.in/ini.v1"
)

type GameConfig struct {
	file *ini.File
	path string
}

// Загрузка INI с сохранением структуры и комментариев
func (c *GameConfig) Load(path string) error {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		PreserveSurroundedQuote:  true, // не трогать кавычки, если появятся
		SpaceBeforeInlineComment: true, // сохранить inline-комментарии
	}, path)
	if err != nil {
		return err
	}
	c.file = cfg
	c.path = path
	return nil
}

// Получить значение
func (c *GameConfig) Get(section, key string) string {
	if c.file == nil {
		return ""
	}
	return c.file.Section(section).Key(key).String()
}

// Обновить значение
func (c *GameConfig) Set(section, key, value string) {
	if c.file == nil {
		return
	}
	c.file.Section(section).Key(key).SetValue(value)
}

// Сохранить обратно в файл
func (c *GameConfig) Save() error {
	if c.file == nil || c.path == "" {
		log.Println("⚠ Save skipped: file or path is nil")
		return nil
	}
	log.Println("💾 Saving INI to:", c.path)
	return c.file.SaveTo(c.path)
}

func (c *GameConfig) Path() string {
	return c.path
}
