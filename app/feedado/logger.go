package feedado

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() (*zap.SugaredLogger, error) {
	var (
		cfg = zap.Config{
			Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
			Development:       false,
			DisableStacktrace: true,
			Encoding:          "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
		err     error
		profile = os.Getenv("APP_PROFILE")
	)

	if profile == "production" {
		cfg.Level.SetLevel(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
