package httpd

import (
	"github.com/donkeywon/gtil/config"
	"time"
)

const (
	DefaultWriteTimeout      = config.Duration(1000 * time.Millisecond)
	DefaultReadTimeout       = config.Duration(1000 * time.Millisecond)
	DefaultReadHeaderTimeout = config.Duration(1000 * time.Millisecond)
	DefaultIdleTimeout       = config.Duration(1000 * time.Millisecond)
	DefaultMonitorInterval   = config.Duration(10 * time.Second)
)

type Config struct {
	Addr              string          `yaml:"addr" mapstructure:"addr" json:"addr"`
	WriteTimeout      config.Duration `yaml:"writeTimeout,omitempty" mapstructure:"writeTimeout,omitempty" json:"writeTimeout,omitempty"`
	ReadTimeout       config.Duration `yaml:"ReadTimeout,omitempty" mapstructure:"ReadTimeout,omitempty" json:"readTimeout,omitempty"`
	ReadHeaderTimeout config.Duration `yaml:"ReadHeaderTimeout,omitempty" mapstructure:"ReadHeaderTimeout,omitempty" json:"readHeaderTimeout,omitempty"`
	IdleTimeout       config.Duration `yaml:"IdleTimeout,omitempty" mapstructure:"IdleTimeout,omitempty" json:"idleTimeout,omitempty"`
	MonitorInterval   config.Duration `yaml:"MonitorInterval,omitempty" mapstructure:"MonitorInterval,omitempty" json:"monitorInterval,omitempty"`
}

func NewConfig(addr string) *Config {
	return &Config{
		Addr:              addr,
		WriteTimeout:      DefaultWriteTimeout,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MonitorInterval:   DefaultMonitorInterval,
	}
}
