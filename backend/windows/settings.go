package windows

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type SettingsWindow struct {
	app    *application.App
	parent *application.WebviewWindow
	child  *application.WebviewWindow
}

// NewSettingsWindow создаёт объект окна настроек и сразу создаёт скрытое дочернее окно
func NewSettingsWindow(app *application.App, parent *application.WebviewWindow) *SettingsWindow {
	sw := &SettingsWindow{
		app:    app,
		parent: parent,
	}

	// Создаём окно один раз, сразу скрытое
	sw.child = sw.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Settings",
		Width:     800,
		Height:    600,
		Frameless: true,
		URL:       "/#/settings", // Хеш-навигация для отдельного контента
		Hidden:    true,
	})

	// Закрытие дочернего окна при закрытии родителя
	sw.parent.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		if sw.child != nil {
			sw.child.Close()
		}
	})

	return sw
}

// Open показывает окно
func (sw *SettingsWindow) Open() {
	if sw.child == nil {
		return
	}
	if sw.child.IsVisible() {
		if !sw.child.IsFocused() {
			sw.child.Focus()
		}
		return
	}
	sw.child.Center()
	sw.child.Show()
}

// Close скрывает окно
func (sw *SettingsWindow) Close() {
	sw.child.Hide()
}
