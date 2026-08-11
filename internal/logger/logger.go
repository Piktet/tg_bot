// Package logger логирование
package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger = zap.NewNop()

// Log возвращает экземпляр логера.
func Log() *zap.Logger {
	return log
}

// InitLogger создание логера.
func InitLogger(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	log = zl
	return nil
}
