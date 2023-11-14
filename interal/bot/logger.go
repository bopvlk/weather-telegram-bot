package bot

import (
	"context"
	"fmt"
	"os"

	"git.foxminded.com.ua/2.4-weather-forecast-bot/interal/models"
	"github.com/sirupsen/logrus"
)

// Logger represents interface for logger
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})

	WarnfCtx(ctx context.Context, format string, args ...interface{})
	DebugfCtx(ctx context.Context, format string, args ...interface{})
	InfofCtx(ctx context.Context, format string, args ...interface{})
	ErrorfCtx(ctx context.Context, format string, args ...interface{})

	DebugCtx(ctx context.Context, args ...interface{})
	InfoCtx(ctx context.Context, args ...interface{})
	ErrorCtx(ctx context.Context, args ...interface{})
}

var entry *logrus.Entry

type LogrusLogger struct {
	logrus *logrus.Logger
	entry  *logrus.Entry
}

func SetUpLogger(cfg *models.Config) (*LogrusLogger, error) {
	level, err := logrus.ParseLevel("trace")
	if err != nil {
		return nil, fmt.Errorf("can't parse log level: %w", err)
	}

	l := &LogrusLogger{
		logrus: logrus.New(),
	}
	l.logrus.SetLevel(level)

	l.logrus.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}

	l.logrus.SetOutput(os.Stdout)
	entry = logrus.NewEntry(l.logrus)
	entry = entry.WithField("service_name", cfg.ServiceName)
	return &LogrusLogger{entry: entry}, nil
}

// Info logs a message at level Info.

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Logf(logrus.InfoLevel, format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Logf(logrus.ErrorLevel, format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Logf(logrus.WarnLevel, format, args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Logf(logrus.DebugLevel, format, args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.entry.Log(logrus.DebugLevel, args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Log(logrus.InfoLevel, args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.entry.Log(logrus.ErrorLevel, args...)
}

func (l *LogrusLogger) InfofCtx(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.InfoLevel, format, args...)
}

func (l *LogrusLogger) DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.DebugLevel, format, args...)
}

func (l *LogrusLogger) WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.WarnLevel, format, args...)
}

func (l *LogrusLogger) ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	l.entry.WithContext(ctx).Logf(logrus.ErrorLevel, format, args...)
}

func (l *LogrusLogger) InfoCtx(ctx context.Context, args ...interface{}) {
	l.entry.WithContext(ctx).Log(logrus.InfoLevel, args...)
}

func (l *LogrusLogger) DebugCtx(ctx context.Context, args ...interface{}) {
	l.entry.WithContext(ctx).Log(logrus.DebugLevel, args...)
}

func (l *LogrusLogger) ErrorCtx(ctx context.Context, args ...interface{}) {
	l.entry.WithContext(ctx).Log(logrus.ErrorLevel, args...)
}
