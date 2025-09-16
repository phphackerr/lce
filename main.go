package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"lce/backend/modules/app_settings"
	"lce/backend/modules/config_editor"
	"lce/backend/modules/config_watcher"
	"lce/backend/modules/i18n"
	"lce/backend/modules/paths_scanner"
	"lce/backend/modules/theming"
	"lce/backend/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "LoD Config Editor",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(i18n.NewI18N()),
			application.NewService(theming.NewThemeService()),
			application.NewService(paths_scanner.NewScanner()),
			application.NewService(config_editor.NewConfigEditor()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "LoD Config Editor",
		Width:     1300,
		Height:    800,
		Frameless: true,
		URL:       "/",
	})

	appSettings := app_settings.NewAppSettings(app)
	app.RegisterService(application.NewService(appSettings))

	settingsWindow := windows.NewSettingsWindow(app, mainWindow)
	app.RegisterService(application.NewService(settingsWindow))

	configWatcher := config_watcher.New(app)
	app.RegisterService(application.NewService(configWatcher))

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
