package logger

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zLogger *zap.SugaredLogger
}

func New(o *config.Options) (*Logger, error) {

	c := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(o.LogLevel)

	if err != nil {
		return nil, err
	}

	c.Level.SetLevel(level)
	c.Encoding = o.LogEncoding
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zl, err := c.Build()

	defer func() { err = zl.Sync() }()

	return &Logger{zl.Sugar()}, err
}

func (l *Logger) Debug(args ...any) {
	l.zLogger.Debug(args)
}

func (l *Logger) Info(args ...any) {
	l.zLogger.Info(args)
}

func (l *Logger) Infoln(args ...any) {
	l.zLogger.Infoln(args)
}

func (l *Logger) Warn(args ...any) {
	l.zLogger.Warn(args)
}

func (l *Logger) Error(args ...any) {
	l.zLogger.Error(args)
}

func (l *Logger) Fatal(args ...any) {
	l.zLogger.Fatal(args)
}
