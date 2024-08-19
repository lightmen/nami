package safe

import (
	"runtime/debug"

	"github.com/lightmen/nami/alog"
)

func Func(f func()) {
	if f == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			alog.Fatal("%v", r)
			alog.Fatal("%s", string(debug.Stack()))
		}
	}()
	f()
}

func Go(f func()) {
	if f == nil {
		return
	}
	go Func(f)
}

func LoopGo(loop func()) {
	if loop == nil {
		return
	}

	go func() {
		for {
			Func(loop)
		}
	}()
}
