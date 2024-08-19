package alog

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
)

type Handler interface {
	//判断是否可以输出日志
	Enabled(context.Context, Level) bool

	Handle(context.Context, Record) error
}

type defaultHandler struct {
	logger    *log.Logger
	calldepth int
}

func newDefaultHandler() *defaultHandler {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	h := &defaultHandler{
		logger:    logger,
		calldepth: 4,
	}

	return h
}

func (*defaultHandler) Enabled(_ context.Context, l Level) bool {
	return l >= LevelInfo
}

func (h *defaultHandler) Handle(ctx context.Context, r Record) error {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	buf.WriteString(r.Level.String())
	buf.WriteByte('\t')
	buf.WriteString(fmt.Sprintf(r.Message, r.args...))

	return h.logger.Output(h.calldepth, buf.String())
}
