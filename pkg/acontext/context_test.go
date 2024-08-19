package acontext

import (
	"context"
	"reflect"
	"testing"
)

func TestWithValue(t *testing.T) {
	key := "key1"
	val1 := "value1"
	val2 := "value2"

	ctx := context.Background()

	actx := New(WithCtx(context.WithValue(context.Background(), key, val1)))

	type args struct {
		parent context.Context
		key    any
		val    any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test1",
			args: args{
				parent: ctx,
				key:    key,
				val:    val1,
			},
			want: val1,
		},
		{
			name: "test2",
			args: args{
				parent: actx,
				key:    key,
				val:    val2,
			},
			want: val2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithValue(tt.args.parent, tt.args.key, tt.args.val)
			val := got.Value(tt.args.key)
			if !reflect.DeepEqual(tt.want, val) {
				t.Errorf("WithValue() = %v, want %v", val, tt.want)
			}
		})
	}

	val := actx.Value(key)
	if val != val2 {
		t.Errorf("WithValue() = %v, want %v", val, val2)
	}
}
