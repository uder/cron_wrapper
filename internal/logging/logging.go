package logging

import "log/slog"

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func getLogger() *slog.Logger {
	logger := slog.Default()
	return logger
}

type SLogger struct {
	logger     *slog.Logger
	enableInfo bool
	enableWarn bool
}

func (l *SLogger) Info(msg string) {
	if l.enableInfo {
		l.logger.Info(msg)
	}
}

func (l *SLogger) Warn(msg string) {
	if l.enableWarn {
		l.logger.Warn(msg)
	}
}

func (l *SLogger) Error(msg string) {
	l.logger.Error(msg)
}

func NewLogger(enableDebug bool) *SLogger {
	return &SLogger{
		logger:     getLogger(),
		enableInfo: enableDebug,
		enableWarn: enableDebug,
	}
}
