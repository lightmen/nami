package nami

import (
	"context"
	"testing"
)

func BenchmarkStruct(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx := context.WithValue(context.Background(), appKey{}, n)

		ctx.Value(appKey{})
	}
}

func BenchmarkString(b *testing.B) {
	key := "test"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx := context.WithValue(context.Background(), key, n)

		ctx.Value(key)
	}
}
