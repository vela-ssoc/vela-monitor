package logger

import (
	"github.com/vela-public/onekit/zapkit"
)

type Logger interface {
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
	Debugf(format string, v ...any)
	Warnf(format string, v ...any)
}

// 让日志只在控制台输出，将 Filename 置空，这样就不会写入文件
var defaultLogger Logger = zapkit.Debug(zapkit.Caller(2, true), zapkit.Console())

func SetLogger(logger Logger) {
	if logger != nil {
		defaultLogger = logger
	}
}

func Infof(format string, v ...any) {
	defaultLogger.Infof(format, v...)
}

func Errorf(format string, v ...any) {
	defaultLogger.Errorf(format, v...)
}

func Debugf(format string, v ...any) {
	defaultLogger.Debugf(format, v...)
}

func Warnf(format string, v ...any) {
	defaultLogger.Warnf(format, v...)
}
