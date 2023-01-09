package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello testHandler"))
}

func TestServer(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	s, err := New(
		Address(":9901"),
		Handler(&testHandler{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	go func() {
		err = s.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	s.Handler.ServeHTTP(w, r)

	bodyStr := w.Body.String()

	t.Logf("rsp: %s", bodyStr)

	assert.Equal(t, "hello testHandler", bodyStr,
		"Then does not order handlers correctly")

	s.Stop(ctx)
}

func TestServer_Endpoint(t *testing.T) {
	addr := "127.0.0.1:9902"
	s, err := New(
		Address(addr),
		Handler(&testHandler{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := s.Endpoint()
	if err != nil {
		t.Fatal(err)
	}

	expect := &url.URL{
		Scheme: s.Name(),
		Host:   addr,
	}

	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("expect: %v, got: %v", expect, got)
	}
}
