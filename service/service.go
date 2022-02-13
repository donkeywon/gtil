package service

import (
	"go.uber.org/zap"
)

type Service interface {
	Name() string
	Open() error
	Close() error
	Shutdown() error
	WithLogger(logger *zap.Logger)
	Statistics() map[string]float64
}
