package httpd

import (
    "context"
    "github.com/donkeywon/gtil/logger"
    "github.com/donkeywon/gtil/service"
    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"
    "testing"
    "time"
)

var (
    cfg    = NewConfig("127.0.0.1:5678")
    log, _ = logger.Default()
    httpd  = New(cfg)
)

func TestHttpD_Open(t *testing.T) {
    httpd.WithLogger(httpd, log)
    httpd.WithContext(context.Background())

    err := service.DoOpen(httpd, context.Background(), log)
    assert.NoError(t, err, "open httpd fail")
}

func TestHttpD_Close(t *testing.T) {
    h := New(cfg)

    err := service.DoClose(h)
    assert.NoError(t, err, "close fail")

    err = service.DoClose(h)
    assert.NoError(t, err, "close fail")
}

func TestHttpD_Ctx(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    l, _ := logger.FromConfig(logger.DefaultConsoleConfig(), zap.Development())
    l.Named("test")

    h := New(cfg)

    err := service.DoOpen(h, ctx, l)
    assert.NoError(t, err, "open httpd fail")

    go func() {
        l.Info("start sleep")
        time.Sleep(time.Second * 5)
        l.Info("cancel")
        cancel()
    }()

    <-h.Closed()
    l.Info("httpd closed")
}
