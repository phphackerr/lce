package utils

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// DialogUtils - это структура, которая будет привязана к фронтенду Wails
type DialogUtils struct {
	ctx context.Context
}

// NewDialogUtils создает новый экземпляр DialogUtils
func NewDialogUtils(ctx context.Context) *DialogUtils {
	return &DialogUtils{
		ctx: ctx,
	}
}

func (d *DialogUtils) SetContext(ctx context.Context) {
	d.ctx = ctx
}

// OpenDirectoryDialog открывает диалог выбора директории
// Возвращает выбранный путь или ошибку.
func (d *DialogUtils) OpenDirectoryDialog(title string) (string, error) {
	selection, err := runtime.OpenDirectoryDialog(d.ctx, runtime.OpenDialogOptions{
		Title:                title,
		CanCreateDirectories: true, // Позволить создание новых директорий
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при открытии диалога выбора директории: %w", err)
	}
	return selection, nil
}
