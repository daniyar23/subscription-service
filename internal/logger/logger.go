package logger

import "go.uber.org/zap"

// New creates a new zap logger instance.
func New() (*zap.Logger, error) {
	return zap.NewProduction()
}
