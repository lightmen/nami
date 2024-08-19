package discovery

import (
	"context"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/lightmen/nami/registry"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestWithTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want *builder
	}{
		{
			name: "test1",
			args: args{
				timeout: time.Duration(123),
			},
			want: &builder{
				timeout: time.Duration(123),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{}
			WithTimeout(tt.args.timeout)(b)
			if !reflect.DeepEqual(b, tt.want) {
				t.Errorf("WithTimeout() = %v, want %v", b, tt.want)
			}
		})
	}
}

type mockDiscovery struct{}

func (m *mockDiscovery) GetService(ctx context.Context, serviceName string) ([]*registry.Instance, error) {
	return nil, nil
}

func (m *mockDiscovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	time.Sleep(time.Microsecond * 500)
	return &testWatch{}, nil
}

type mockConn struct{}

func (m *mockConn) UpdateState(resolver.State) error {
	return nil
}

func (m *mockConn) ReportError(error) {}

func (m *mockConn) NewAddress(addresses []resolver.Address) {}

func (m *mockConn) NewServiceConfig(serviceConfig string) {}

func (m *mockConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}

func Test_builder_Build(t *testing.T) {
	b := NewBuilder(&mockDiscovery{})

	_, err := b.Build(
		resolver.Target{
			URL: url.URL{
				Scheme: resolver.GetDefaultScheme(),
				Path:   "grpc://authority/endpoint",
			},
		},
		&mockConn{},
		resolver.BuildOptions{},
	)
	if err != nil {
		t.Errorf("expect no error, get: %v", err)
		return
	}

	timeoutBuilder := NewBuilder(&mockDiscovery{}, WithTimeout(0))
	_, err = timeoutBuilder.Build(
		resolver.Target{
			URL: url.URL{
				Scheme: resolver.GetDefaultScheme(),
				Path:   "grpc://authority/endpoint",
			},
		},
		&mockConn{},
		resolver.BuildOptions{},
	)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}
