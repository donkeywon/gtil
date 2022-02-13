package httpd

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net/http"
)

const (
	Name = "httpd"
)

type HttpD struct {
	config *Config

	server *http.Server

	closed chan struct{}

	ctx context.Context

	err error
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
			s.err = multierr.Append(s.err, err)
		}
	}()

	return nil
}

func (s *HttpD) Close() error {
	select {
	case <-s.Closed():
		return nil
	default:
		defer close(s.closed)
		err := s.server.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HttpD) Shutdown() error {
	select {
	case <-s.Closed():
		return nil
	default:
		defer close(s.closed)
		err := s.server.Shutdown(s.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HttpD) Closed() <-chan struct{} {
	return s.closed
}

func (s *HttpD) WithLogger(logger *zap.Logger) {

}

func (s *HttpD) SetHandler(router *mux.Router) {
	s.server.Handler = router
}

func (s *HttpD) Serve() {
	err := s.server.ListenAndServe()
	s.err = multierr.Append(s.err, err)
}

func (s *HttpD) LastError() error {
	return s.err
}

func (s *HttpD) Statistics() map[string]float64 {
	return nil
}
