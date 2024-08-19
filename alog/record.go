package alog

import (
	"time"

	"github.com/lightmen/nami/pkg/cast"
	"google.golang.org/protobuf/proto"
)

type Record struct {
	Time time.Time

	Message string

	Level Level

	info callerInfo

	args []any
}

func NewRecord(t time.Time, level Level, msg string, info callerInfo) Record {
	return Record{
		Time:    t,
		Message: msg,
		Level:   level,
		info:    info,
	}
}

func (r *Record) BindArgs(args ...any) {
	argFmt := []any{}
	for _, v := range args {
		if _, ok := v.(proto.Message); ok {
			argFmt = append(argFmt, cast.ToJson(v))
		} else {
			argFmt = append(argFmt, v)
		}
	}
	r.args = argFmt
}
