package log

import (
	"log"
	"os"
)

type defaultLogger struct {
	*log.Logger
}

func Default() Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	d := &defaultLogger{
		Logger: logger,
	}

	return d
}

func (d *defaultLogger) Debug(format string, args ...interface{}) {
	d.Printf(format, args...)
}

func (d *defaultLogger) Info(format string, args ...interface{}) {
	d.Printf(format, args...)
}

func (d *defaultLogger) Error(format string, args ...interface{}) {
	d.Printf(format, args...)
}

func (d *defaultLogger) Fatal(format string, args ...interface{}) {
	d.Printf(format, args...)
}
