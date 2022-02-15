package httpd

import (
    "github.com/donkeywon/gtil/service"
    "github.com/gorilla/mux"
    "net/http"
)

const (
    Name = "httpd"
)

type HttpD struct {
    *service.BaseService

    config *Config

    server *http.Server
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

func New(config *Config) *HttpD {
    return &HttpD{
        BaseService: service.NewBase(),
        config:      config,
        server:      newHTTPServer(config),
    }
}

func (s *HttpD) Name() string {
    return Name
}

func (s *HttpD) Open() error {
    go s.serve()
    return nil
}

func (s *HttpD) Close() error {
    return s.server.Close()
}

func (s *HttpD) Shutdown() error {
    return s.server.Shutdown(s.Context())
}

func (s *HttpD) SetHandler(router *mux.Router) {
    s.server.Handler = router
}

func (s *HttpD) serve() {
    s.AppendError(s.server.ListenAndServe())
}
