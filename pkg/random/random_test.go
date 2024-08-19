package random

import (
	"testing"
	"unsafe"
)

func TestString(t *testing.T) {
	s1 := "abcdefg"
	bts := []byte(s1)
	s2 := unsafe.String(&bts[0], len(bts))
	if s1 != s2 {
		t.Fatal("String() error")
	}

	if len(String(16)) != 16 {
		t.Fatal("String() len is wrong")
	}
	if len(Bytes(16)) != 16 {
		t.Fatal("Bytes() len is wrong")
	}
}

func BenchmarkString(b *testing.B) {
	s := "abcd0123"
	bts := []byte(s)
	ln := len(bts)

	for i := 0; i < b.N; i++ {
		_ = unsafe.String(&bts[0], ln)
	}
}

func TestIsRanded(t *testing.T) {
	type args struct {
		weight int32
		opts   []Option
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				weight: 1,
				opts: []Option{
					WithRandVal(1),
				},
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				weight: 0,
				opts:   []Option{},
			},
			want: false,
		},
		{
			name: "test3",
			args: args{
				weight: 10000,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRanded(tt.args.weight, tt.args.opts...); got != tt.want {
				t.Errorf("IsRanded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArray(t *testing.T) {
	iarr := []int{1, 2, 4, 5, 6}

	t.Logf("%v", Array(iarr))

}
