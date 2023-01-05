package log

var gLogger = Default()

func SetLogger(logger Logger) {
	gLogger = logger
}

func Info(format string, args ...interface{}) {
	gLogger.Info(format, args...)
}

func Error(format string, args ...interface{}) {
	gLogger.Error(format, args...)
}
