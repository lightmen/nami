package discovery

import (
	"context"
	"testing"
	"time"

	"github.com/lightmen/nami/registry"
	"google.golang.org/grpc/resolver"
)

type testWatch struct {
	err error

	count uint
}

func (m *testWatch) Next() ([]*registry.Instance, error) {
	time.Sleep(time.Millisecond * 200)
	if m.count > 1 {
		return nil, nil
	}
	m.count++
	ins := []*registry.Instance{
		{
			ID:        "mock_ID",
			Name:      "mock_Name",
			Endpoints: []string{"grpc://127.0.0.1"},
		},
		{
			ID:        "mock_ID2",
			Name:      "mock_Name2",
			Endpoints: []string{""},
		},
	}
	return ins, m.err
}

// Watch creates a watcher according to the service name.
func (m *testWatch) Stop() error {
	return m.err
}

type testClientConn struct {
	resolver.ClientConn // For unimplemented functions
	te                  *testing.T
}

func (t *testClientConn) UpdateState(s resolver.State) error {
	t.te.Log("UpdateState", s)
	return nil
}

func TestWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &discoveryResolver{
		w:      &testWatch{},
		cc:     &testClientConn{te: t},
		ctx:    ctx,
		cancel: cancel,
	}

	r.ResolveNow(resolver.ResolveNowOptions{})
	go func() {
		time.Sleep(time.Second * 2)
		r.Close()
	}()
	r.watch()
	t.Log("watch goroutine exited after 2 second")
}
