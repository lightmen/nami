package alog

import (
	"testing"
)

func TestCallerInfo(t *testing.T) {
	info := CallerInfo(0)
	if info.funcName != "alog.TestCallerInfo" {
		t.Errorf("test failed, got: %s", info.funcName)
		return
	}

	t.Logf("info: %v", info)
}
