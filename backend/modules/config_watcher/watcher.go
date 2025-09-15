package config_watcher

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type ConfigWatcher struct {
	app      *application.App
	filePath string
	watcher  *fsnotify.Watcher
	stop     chan struct{}
}

// New создаёт вотчер
func New(app *application.App) *ConfigWatcher {
	return &ConfigWatcher{
		app:  app,
		stop: make(chan struct{}),
	}
}

// StartWatching — вызывается из фронтенда
func (cw *ConfigWatcher) StartWatching(path string) error {
	cw.filePath = path

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	cw.watcher = w

	dir := filepath.Dir(cw.filePath)
	if err := w.Add(dir); err != nil {
		return err
	}

	go cw.run()
	return nil
}

func (cw *ConfigWatcher) run() {
	for {
		select {
		case event := <-cw.watcher.Events:
			if filepath.Clean(event.Name) == filepath.Clean(cw.filePath) &&
				(event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create) {

				log.Println("⚡ Config file changed:", event)

				// эмитим событие для фронтенда
				cw.app.Event.Emit("config-changed", cw.filePath)
			}

		case err := <-cw.watcher.Errors:
			log.Println("Watcher error:", err)

		case <-cw.stop:
			_ = cw.watcher.Close()
			return
		}
	}
}

// StopWatching — останавливает наблюдение
func (cw *ConfigWatcher) StopWatching() {
	if cw.stop != nil {
		close(cw.stop)
	}
}
