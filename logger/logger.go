package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	DefaultTimeFormatLayout  = "2006-01-02 15:04:05.000000"
	DefaultLevel             = zap.DebugLevel
	DefaultIsDev             = false
	DefaultDisableCaller     = false
	DefaultDisableStacktrace = false

	DefaultEncoding                = "console"
	DefaultEncoderMessageKey       = "msg"
	DefaultEncoderLevelKey         = "lvl"
	DefaultEncoderNameKey          = "logger"
	DefaultEncoderTimeKey          = "ts"
	DefaultEncoderCallerKey        = "caller"
	DefaultEncoderFunctionKey      = "fn"
	DefaultEncoderStacktraceKey    = "stacktrace"
	DefaultEncoderSkipLineEncoding = false
	DefaultEncoderLineEnding       = "\n"
	DefaultEncoderLevelEncoder     = "lowercase"
	DefaultEncoderTimeEncoder      = "layout"
	DefaultEncoderDurationEncoder  = "string"
	DefaultEncoderCallerEncoder    = "short"
	DefaultEncoderNameEncoder      = "full"
	DefaultEncoderConsoleSeparator = "\t"
)

var (
	DefaultOutputPath      = []string{"stdout"}
	DefaultErrorOutputPath = []string{"stderr"}
)

func defaultEncoderConfig() zapcore.EncoderConfig {
	var le zapcore.LevelEncoder
	_ = le.UnmarshalText([]byte(DefaultEncoderLevelEncoder))

	var te zapcore.TimeEncoder
	if DefaultEncoderTimeEncoder == "layout" {
		te = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(ts.Format(DefaultTimeFormatLayout))
		}
	} else {
		_ = te.UnmarshalText([]byte(DefaultEncoderTimeEncoder))
	}

	var de zapcore.DurationEncoder
	_ = de.UnmarshalText([]byte(DefaultEncoderDurationEncoder))

	var ce zapcore.CallerEncoder
	_ = ce.UnmarshalText([]byte(DefaultEncoderCallerEncoder))

	var ne zapcore.NameEncoder
	_ = ne.UnmarshalText([]byte(DefaultEncoderNameEncoder))

	config := zapcore.EncoderConfig{
		MessageKey:       DefaultEncoderMessageKey,
		LevelKey:         DefaultEncoderLevelKey,
		TimeKey:          DefaultEncoderTimeKey,
		NameKey:          DefaultEncoderNameKey,
		CallerKey:        DefaultEncoderCallerKey,
		FunctionKey:      DefaultEncoderFunctionKey,
		StacktraceKey:    DefaultEncoderStacktraceKey,
		SkipLineEnding:   DefaultEncoderSkipLineEncoding,
		LineEnding:       DefaultEncoderLineEnding,
		EncodeLevel:      le,
		EncodeTime:       te,
		EncodeDuration:   de,
		EncodeCaller:     ce,
		EncodeName:       ne,
		ConsoleSeparator: DefaultEncoderConsoleSeparator,
	}

	return config
}

func DefaultConfig() zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(DefaultLevel),
		Development:       DefaultIsDev,
		DisableCaller:     DefaultDisableCaller,
		DisableStacktrace: DefaultDisableStacktrace,
		Sampling:          nil,
		Encoding:          DefaultEncoding,
		EncoderConfig:     defaultEncoderConfig(),
		OutputPaths:       DefaultOutputPath,
		ErrorOutputPaths:  DefaultErrorOutputPath,
		InitialFields:     nil,
	}
}

func Default(option ...zap.Option) (*zap.Logger, error) {
	return DefaultConfig().Build(option...)
}

func FromConfig(config zap.Config, option ...zap.Option) (*zap.Logger, error) {
	return config.Build(option...)
}
