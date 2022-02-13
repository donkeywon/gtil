package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestDefault(t *testing.T) {
	l, _ := Default()
	l.Info("WTF", zap.String("tag", "value"))
}

func TestParseConfig(t *testing.T) {
	f, err := ioutil.ReadFile("D:/workspace/code/go/gtil/logger/log.yaml")
	assert.NoError(t, err, "read config fail")

	config := &zap.Config{}
	err = yaml.Unmarshal(f, config)
	assert.NoError(t, err, "unmarshal fail")

	l, err := config.Build()
	assert.NoError(t, err, "build logger fail")
	l.Info("WTF", zap.String("tag", "value"))
	l.Info("QWE", zap.String("t", "v"))
}
