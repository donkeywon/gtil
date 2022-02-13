package httpc

import (
	"github.com/donkeywon/gtil/config"
	"time"
)

const (
	DefaultTimeout                = config.Duration(time.Second * 5)
	DefaultDisableKeepAlives      = false
	DefaultDisableCompression     = false
	DefaultMaxIdleConns           = 0
	DefaultMaxIdleConnsPerHost    = 2
	DefaultMaxConnsPerHost        = 0
	DefaultIdleConnTimeout        = config.Duration(0)
	DefaultResponseHeaderTimeout  = config.Duration(0)
	DefaultExpectContinueTimeout  = config.Duration(0)
	DefaultMaxResponseHeaderBytes = 0
	DefaultWriteBufferSize        = 4096
	DefaultReadBufferSize         = 4096
)

type Config struct {
	Timeout                config.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	DisableKeepAlives      bool            `yaml:"disableKeepAlives,omitempty" json:"disableKeepAlives,omitempty"`
	DisableCompression     bool            `yaml:"disableCompression,omitempty" json:"disableCompression,omitempty"`
	MaxIdleConns           int             `yaml:"maxIdleConns,omitempty" json:"maxIdleConns,omitempty"`
	MaxIdleConnsPerHost    int             `yaml:"maxIdleConnsPerHost,omitempty" json:"maxIdleConnsPerHost,omitempty"`
	MaxConnsPerHost        int             `yaml:"maxConnsPerHost,omitempty" json:"maxConnsPerHost,omitempty"`
	IdleConnTimeout        config.Duration `yaml:"idleConnTimeout,omitempty" json:"idleConnTimeout,omitempty"`
	ResponseHeaderTimeout  config.Duration `yaml:"responseHeaderTimeout,omitempty" json:"responseHeaderTimeout,omitempty"`
	ExpectContinueTimeout  config.Duration `yaml:"expectContinueTimeout,omitempty" json:"expectContinueTimeout,omitempty"`
	MaxResponseHeaderBytes int64           `yaml:"maxResponseHeaderBytes,omitempty" json:"maxResponseHeaderBytes,omitempty"`
	WriteBufferSize        int             `yaml:"writeBufferSize,omitempty" json:"writeBufferSize,omitempty"`
	ReadBufferSize         int             `yaml:"readBufferSize,omitempty" json:"readBufferSize,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Timeout:                DefaultTimeout,
		DisableKeepAlives:      DefaultDisableKeepAlives,
		DisableCompression:     DefaultDisableCompression,
		MaxIdleConns:           DefaultMaxIdleConns,
		MaxIdleConnsPerHost:    DefaultMaxIdleConnsPerHost,
		MaxConnsPerHost:        DefaultMaxConnsPerHost,
		IdleConnTimeout:        DefaultIdleConnTimeout,
		ResponseHeaderTimeout:  DefaultResponseHeaderTimeout,
		ExpectContinueTimeout:  DefaultExpectContinueTimeout,
		MaxResponseHeaderBytes: DefaultMaxResponseHeaderBytes,
		WriteBufferSize:        DefaultWriteBufferSize,
		ReadBufferSize:         DefaultReadBufferSize,
	}
}
