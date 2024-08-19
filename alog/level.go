package alog

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}

	return "UNKNOWN"
}
