package copyutil

import (
	"reflect"
	"testing"
)

type Item struct {
	Type     int32
	ID       int32
	DeltaNum int64
}

func TestClone(t *testing.T) {
	src := []*Item{
		{Type: 1, ID: 2, DeltaNum: 500000},
		{Type: 2, ID: 3, DeltaNum: 100},
	}

	res := Clone(src)
	if res == nil {
		t.Fatal("res is nil")
	}
	dst, ok := res.([]*Item)
	if !ok {
		t.Fatal("clone type is wrong")
	}
	if !reflect.DeepEqual(src, dst) {
		t.Fatal("clone data is not equal")
	}
}

func BenchmarkClone(b *testing.B) {
	src := []*Item{
		{Type: 1, ID: 2, DeltaNum: 500000},
		{Type: 2, ID: 3, DeltaNum: 100},
	}

	for i := 0; i < b.N; i++ {
		_ = Clone(src)
	}
}
