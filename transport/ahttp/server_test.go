package ahttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello testHandler"))
}

func TestServer(t *testing.T) {
	s, err := New(":9901")
	if err != nil {
		t.Fatal(err)
	}

	s.Handle("/", &testHandler{})

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	ctx := context.Background()
	go func() {
		err = s.Start(ctx)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()

	s.Handler.ServeHTTP(w, r)

	bodyStr := w.Body.String()

	t.Logf("rsp: %s", bodyStr)

	assert.Equal(t, "hello testHandler", bodyStr,
		"Then does not order handlers correctly")

	s.Stop(ctx)
}

func TestServer_Endpoint(t *testing.T) {
	addr := "127.0.0.1:9902"
	s, err := New(addr)
	if err != nil {
		t.Fatal(err)
	}

	handler := &testHandler{}
	s.HandleFunc("/", handler.ServeHTTP)

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
