package logger

import (
	"fmt"
)

func Debug(msg string, fields ...any) {
	Logger.Debug(fmt.Sprintf(msg, fields...))
}

func Info(msg string, fields ...any) {
	Logger.Info(fmt.Sprintf(msg, fields...))
}

func Warn(msg string, fields ...any) {
	Logger.Warn(fmt.Sprintf(msg, fields...))
}

func Error(msg string, fields ...any) {
	Logger.Error(fmt.Sprintf(msg, fields...))
}

func Fatal(msg string, fields ...any) {
	Logger.Fatal(fmt.Sprintf(msg, fields...))
}
