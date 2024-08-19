package alog

import (
	"runtime"
	"strings"
)

type callerInfo struct {
	funcName string
	fileName string
	line     int
}

func CallerInfo(skip int) callerInfo {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return callerInfo{}
	}

	funcName := runtime.FuncForPC(pc).Name()
	nameIdx := strings.LastIndex(funcName, "/")
	funcName = funcName[nameIdx+1:]

	return callerInfo{
		funcName: funcName,
		fileName: file,
		line:     line,
	}
}
