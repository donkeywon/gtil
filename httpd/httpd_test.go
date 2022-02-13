package httpd

import (
	"context"
	"github.com/donkeywon/gtil/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	cfg    = NewConfig("127.0.0.1:5678")
	log, _ = logger.Default()
	httpd  = New(cfg, context.Background())
)

func TestHTTPd_Open(t *testing.T) {
	httpd.WithLogger(log)

	err := httpd.Open()
	assert.NoError(t, err, "open httpd fail")
}

func TestHTTPd_Close(t *testing.T) {
	h := New(cfg, context.Background())
	h.WithLogger(log)

	err := h.Close()
	assert.NoError(t, err, "close fail")

	err = h.Close()
	assert.NoError(t, err, "close fail")
}
