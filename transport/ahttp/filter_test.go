package ahttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tagFilterFunc(tag string) FilterFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

func TestFilterChain(t *testing.T) {
	f1 := tagFilterFunc("f1")
	f2 := tagFilterFunc("f2")

	filters := FilterChain(f1, f2)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("next"))
	})
	handler := filters(next)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(w, r)

	// t.Logf("result: %v", w.Body.String())
	assert.Equal(t, "f1f2next", w.Body.String(),
		"Then does not order handlers correctly")
}
