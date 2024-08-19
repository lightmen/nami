package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestAppendToClientContext(t *testing.T) {
	type args struct {
		md Metadata
		kv []string
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "kratos",
			args: args{Metadata{}, []string{"hello", "kratos", "env", "dev"}},
			want: Metadata{"hello": "kratos", "env": "dev"},
		},
		{
			name: "hello",
			args: args{Metadata{"hi": "https://go-kratos.dev/"}, []string{"hello", "kratos", "env", "dev"}},
			want: Metadata{"hello": "kratos", "env": "dev", "hi": "https://go-kratos.dev/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewClientContext(context.Background(), tt.args.md)
			ctx = AppendToClientContext(ctx, tt.args.kv...)
			md, ok := FromClientContext(ctx)
			if !ok {
				t.Errorf("FromServerContext() = %v, want %v", ok, true)
			}
			if !reflect.DeepEqual(md, tt.want) {
				t.Errorf("metadata = %v, want %v", md, tt.want)
			}
		})
	}
}

func TestMergeToClientContext(t *testing.T) {
	ctx := context.Background()
	md := Metadata{"key1": "value1", "key2": "value2"}
	ctx = MergeToClientContext(ctx, md)
	newMd, ok := FromClientContext(ctx)
	if !ok {
		t.Error("failed to get metadata from context")
	}
	if newMd.Get("key1") != "value1" {
		t.Errorf("expected key1=value1, got %s = %s", newMd.Get("key1"), newMd.Get("key2"))
	}
	if newMd.Get("key2") != "value2" {
		t.Errorf("expected key2=value2, got %s = %s", newMd.Get("key2"), newMd.Get("key1"))
	}
	if len(newMd) != 2 {
		t.Errorf("expected metadata of length 2, got %d", len(newMd))
	}
}

func TestMergeToClientContextEmptyMetadata(t *testing.T) {
	ctx := context.Background()
	md := Metadata{}
	ctx = MergeToClientContext(ctx, md)
	newMd, ok := FromClientContext(ctx)
	if !ok {
		t.Error("failed to get metadata from context")
	}
	if len(newMd) != 0 {
		t.Errorf("expected empty metadata, got %d metadata items", len(newMd))
	}
}

func TestMergeToClientContextNilContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	md := Metadata{"key1": "value1", "key2": "value2"}
	MergeToClientContext(nil, md)
}
