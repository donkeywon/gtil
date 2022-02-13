package httpc

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

const (
	Name = "httpc"
)

type HttpC struct {
	*http.Client

	config *Config
	logger *zap.Logger

	closed chan struct{}
	ctx    context.Context
}

func New(config *Config, ctx context.Context) *HttpC {
	trans := &http.Transport{
		DisableKeepAlives:      config.DisableKeepAlives,
		DisableCompression:     config.DisableCompression,
		MaxIdleConns:           config.MaxIdleConns,
		MaxIdleConnsPerHost:    config.MaxIdleConnsPerHost,
		IdleConnTimeout:        config.IdleConnTimeout.ToDuration(),
		ResponseHeaderTimeout:  config.ResponseHeaderTimeout.ToDuration(),
		ExpectContinueTimeout:  config.ExpectContinueTimeout.ToDuration(),
		MaxResponseHeaderBytes: config.MaxResponseHeaderBytes,
		WriteBufferSize:        config.WriteBufferSize,
		ReadBufferSize:         config.ReadBufferSize,
	}

	c := &http.Client{
		Timeout:   config.Timeout.ToDuration(),
		Transport: trans,
	}

	return &HttpC{
		Client: c,
		config: config,
		closed: make(chan struct{}),
		ctx:    ctx,
	}
}

func (h *HttpC) Name() string {
	return Name
}

func (h *HttpC) Open() error {
	return nil
}

func (h *HttpC) Close() error {
	return nil
}

func (h *HttpC) Shutdown() error {
	return nil
}

func (h *HttpC) WithLogger(logger *zap.Logger) {
	h.logger = logger.Named(h.Name())
}

func (h *HttpC) Statistics() map[string]float64 {
	return nil
}
