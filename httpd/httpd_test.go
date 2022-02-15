package httpd

import (
	"context"
	"github.com/donkeywon/gtil/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

var (
	cfg    = NewConfig("127.0.0.1:5678")
	log, _ = logger.Default()
	httpd  = New(cfg, context.Background())
)

func TestHttpD_Open(t *testing.T) {
	httpd.WithLogger(log)

	err := httpd.Open()
	assert.NoError(t, err, "open httpd fail")
}

func TestHttpD_Close(t *testing.T) {
	h := New(cfg, context.Background())
	h.WithLogger(log)

	err := h.Close()
	assert.NoError(t, err, "close fail")

	err = h.Close()
	assert.NoError(t, err, "close fail")
}

func TestHttpD_Ctx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	lc := logger.DefaultConsoleConfig()
	l, _ := logger.FromConfig(lc, zap.Development())
	l.Named("test")

	h := New(cfg, ctx)
	h.WithLogger(l)

	err := h.Open()
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
