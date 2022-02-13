package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	JsonEncoding    = "json"
	ConsoleEncoding = "console"

	LevelEncoderCapital      = "capital"
	LevelEncoderCapitalColor = "capitalColor"
	LevelEncoderColor        = "color"
	LevelEncoderLowercase    = "lowercase"

	TimeEncoderRFC3339Nano = "rfc3339nano"
	TimeEncoderRFC3339     = "rfc3339"
	TimeEncoderISO8601     = "iso8601"
	TimeEncoderMillis      = "millis"
	TimeEncoderNanos       = "nanos"
	TimeEncoderSecond      = "second"
	TimeEncoderLayout      = "2006-01-02 15:04:05.000000"

	DurationEncoderString = "string"
	DurationEncoderNanos  = "nanos"
	DurationEncoderMillis = "ms"
	DurationEncoderSecond = "second"

	CallerEncoderFull  = "full"
	CallerEncoderShort = "short"

	NameEncoderFull = "full"

	DefaultLevel             = zap.DebugLevel
	DefaultIsDev             = false
	DefaultDisableCaller     = false
	DefaultDisableStacktrace = false

	DefaultEncoding                = JsonEncoding
	DefaultEncoderMessageKey       = "msg"
	DefaultEncoderLevelKey         = "lvl"
	DefaultEncoderNameKey          = "logger"
	DefaultEncoderTimeKey          = "ts"
	DefaultEncoderCallerKey        = "caller"
	DefaultEncoderFunctionKey      = "fn"
	DefaultEncoderStacktraceKey    = "stacktrace"
	DefaultEncoderSkipLineEncoding = false
	DefaultEncoderLineEnding       = "\n"
	DefaultEncoderLevelEncoder     = LevelEncoderLowercase
	DefaultEncoderTimeEncoder      = TimeEncoderLayout
	DefaultEncoderDurationEncoder  = DurationEncoderString
	DefaultEncoderCallerEncoder    = CallerEncoderShort
	DefaultEncoderNameEncoder      = NameEncoderFull
	DefaultEncoderConsoleSeparator = "\t"
)

var (
	DefaultOutputPath      = []string{"stdout"}
	DefaultErrorOutputPath = []string{"stderr"}
)

func buildTimeEncoder(enc string) zapcore.TimeEncoder {
	var te zapcore.TimeEncoder
	if enc[0] >= '0' && enc[1] <= '9' {
		te = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(ts.Format(enc))
		}
	} else {
		_ = te.UnmarshalText([]byte(enc))
	}

	return te
}

func defaultEncoderConfig() *zapcore.EncoderConfig {
	var le zapcore.LevelEncoder
	_ = le.UnmarshalText([]byte(DefaultEncoderLevelEncoder))

	var te zapcore.TimeEncoder
	te = buildTimeEncoder(DefaultEncoderTimeEncoder)

	var de zapcore.DurationEncoder
	_ = de.UnmarshalText([]byte(DefaultEncoderDurationEncoder))

	var ce zapcore.CallerEncoder
	_ = ce.UnmarshalText([]byte(DefaultEncoderCallerEncoder))

	var ne zapcore.NameEncoder
	_ = ne.UnmarshalText([]byte(DefaultEncoderNameEncoder))

	config := &zapcore.EncoderConfig{
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

func DefaultConfig() *zap.Config {
	return &zap.Config{
		Level:             zap.NewAtomicLevelAt(DefaultLevel),
		Development:       DefaultIsDev,
		DisableCaller:     DefaultDisableCaller,
		DisableStacktrace: DefaultDisableStacktrace,
		Sampling:          nil,
		Encoding:          DefaultEncoding,
		EncoderConfig:     *defaultEncoderConfig(),
		OutputPaths:       DefaultOutputPath,
		ErrorOutputPaths:  DefaultErrorOutputPath,
		InitialFields:     nil,
	}
}

func DefaultConsoleConfig() *zap.Config {
	c := DefaultConfig()
	c.Encoding = ConsoleEncoding
	return c
}

func DefaultJsonConfig() *zap.Config {
	c := DefaultConfig()
	c.Encoding = JsonEncoding
	return c
}

func Default(option ...zap.Option) (*zap.Logger, error) {
	return DefaultConfig().Build(option...)
}

func FromConfig(config zap.Config, option ...zap.Option) (*zap.Logger, error) {
	return config.Build(option...)
}
