package etcd

import (
	"reflect"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestNew(t *testing.T) {
	type args struct {
		client *clientv3.Client
		opts   []Option
	}
	tests := []struct {
		name string
		args args
		want *Registry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.client, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
