package recovery

import (
	"context"
	"testing"
)

func TestRecovery(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Error("fail")
		}
	}()

	next := func(ctx context.Context, req any) (any, error) {
		panic("panic reason")
	}
	_, e := Recovery()(next)(context.Background(), "panic")
	t.Logf("succ and reason is %v", e)
}
