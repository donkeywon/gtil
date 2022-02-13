package httpd

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

const (
	Name = "httpd"
)

type HttpD struct {
	config *Config

	logger *zap.Logger

	server *http.Server

	closed chan struct{}

	ctx context.Context
}

func newHTTPServer(config *Config) *http.Server {
	return &http.Server{
		Addr:              config.Addr,
		ReadTimeout:       config.ReadTimeout.ToDuration(),
		ReadHeaderTimeout: config.ReadHeaderTimeout.ToDuration(),
		WriteTimeout:      config.WriteTimeout.ToDuration(),
		IdleTimeout:       config.IdleTimeout.ToDuration(),
	}
}

func New(config *Config, ctx context.Context) *HttpD {
	return &HttpD{
		ctx:    ctx,
		config: config,
		server: newHTTPServer(config),
		closed: make(chan struct{}),
	}
}

func (s *HttpD) Name() string {
	return Name
}

func (s *HttpD) Open() error {
	go s.Serve()

	go func() {
		<-s.ctx.Done()
		err := s.Close()
		if err != nil {
			s.logger.Error("Close http server fail", zap.Error(err))
		}
	}()

	s.logger.Info("Open http Server")

	return nil
}

func (s *HttpD) Close() error {
	select {
	case <-s.Closed():
		return nil
	default:
		defer close(s.closed)
		s.logger.Info("Closing")

		err := s.server.Close()
		if err != nil {
			return err
		}

		s.logger.Info("Closed")
	}

	return nil
}

func (s *HttpD) Shutdown() error {
	select {
	case <-s.Closed():
		return nil
	default:
		defer close(s.closed)
		s.logger.Info("Shutdown start")

		err := s.server.Shutdown(s.ctx)
		if err != nil {
			return err
		}

		s.logger.Info("Shutdown done")
	}
	return nil
}

func (s *HttpD) Closed() <-chan struct{} {
	return s.closed
}

func (s *HttpD) WithLogger(logger *zap.Logger) {
	s.logger = logger.Named(s.Name())
}

func (s *HttpD) SetHandler(router *mux.Router) {
	s.server.Handler = router
}

func (s *HttpD) Serve() {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("Serve Fail", zap.Error(err))
	}
}
