package alog

import (
	"testing"
)

func TestInfoCtx(t *testing.T) {
	// SetDefault(New(NewFileHandler("test", ".")))
	a := 324
	b := "lily"
	InfoCtx(nil, "just test: %d, name: %s", a, b)
}
