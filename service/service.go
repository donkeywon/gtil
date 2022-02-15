package service

import (
	"context"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Service interface {
	Name() string
	Open() error
	Close() error
	Shutdown() error
	WithLogger(logger *zap.Logger)
	Statistics() map[string]float64
	AppendError(err error)
	LastError() error
}

type BaseService struct {
	logger *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
	err    error
}

func (bs *BaseService) AppendError(err error) {
	bs.err = multierr.Append(bs.err, err)
}

func (bs *BaseService) LastError() error {
	return bs.err
}

func (bs *BaseService) WithLogger(logger *zap.Logger) {
	bs.logger = logger.Named(bs.Name())
}

func (bs *BaseService) Name() string {
	panic("Override me")
}
