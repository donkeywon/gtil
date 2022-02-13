package sink

import (
	"github.com/donkeywon/gtil/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestHttp_Write(t *testing.T) {
	lc := logger.DefaultConfig()
	lc.Encoding = "json"
	lc.OutputPaths = append(lc.OutputPaths, "http://httpbin.org/post")

	l, err := logger.FromConfig(lc)
	assert.NoError(t, err, "init logger fail")

	l.Info("WTF message", zap.String("tag", "value"))
}
