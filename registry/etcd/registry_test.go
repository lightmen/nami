package etcd

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/lightmen/nami/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func TestRegistry(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	instance := &registry.Instance{
		ID:        "0",
		Name:      "helloworld",
		Endpoints: []string{"http://127.0.0.1:8000", "grpc://127.0.0.1:9000"},
	}

	r := New(client)

	err = r.Watch(ctx, instance.Name)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(ctx, instance)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(ctx, instance.Name)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 && !reflect.DeepEqual(res, instance) {
		t.Errorf("instant inspect %+v, got: %+v", instance, res)
	}

	time.Sleep(time.Second * 1)

	t.Logf("start unregister ectd")

	err = r.Unregister(ctx, instance)
	if err != nil {
		t.Fatal(err)
	}

	res, err = r.GetService(ctx, instance.Name)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 0 {
		t.Errorf("res is not releaseed")
	}
}
