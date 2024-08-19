package chain

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tagMiddleware(tag string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

func funcsEqual(f1, f2 any) bool {
	val1 := reflect.ValueOf(f1)
	val2 := reflect.ValueOf(f2)

	return val1.Pointer() == val2.Pointer()
}

func TestNew(t *testing.T) {
	type args struct {
		middlewares []Middleware
	}

	m1 := tagMiddleware("tag1")
	m2 := tagMiddleware("tag2")
	mids := []Middleware{m1, m2}
	tests := []struct {
		name        string
		args        args
		middlewares []Middleware
	}{
		{
			name: "test1",
			args: args{
				middlewares: []Middleware{m1, m2},
			},
			middlewares: mids,
		},
		{
			name: "test2",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.middlewares...)
			gotChain, ok := got.(*chain)
			if !ok {
				t.Errorf("New() is not *chain")
			}

			if len(gotChain.middlewares) != len(tt.middlewares) {
				t.Errorf("got middlewares %d, want: %d", len(gotChain.middlewares), len(tt.middlewares))
			}

			for idx, c := range gotChain.middlewares {
				assert.True(t, funcsEqual(c, tt.middlewares[idx]),
					"New does not add constructors correctly")
			}

		})
	}
}

var testApp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app"))
})

func Test_chain_Then(t *testing.T) {
	t1 := tagMiddleware("t1")
	t2 := tagMiddleware("t2")
	t3 := tagMiddleware("t3")

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	c := New(t1, t2, t3).Then(testApp)

	c.ServeHTTP(w, r)

	assert.Equal(t, "t1t2t3app", w.Body.String(),
		"Then does not order handlers correctly")
}

func Test_chain_Append(t *testing.T) {
	t1 := tagMiddleware("t1")
	t2 := tagMiddleware("t2")
	t3 := tagMiddleware("t3")

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	chain := New(t1)

	chain.Append(t2, t3)
	c := chain.Then(testApp)

	c.ServeHTTP(w, r)

	assert.Equal(t, "t1t2t3app", w.Body.String(),
		"Then does not order handlers correctly")
}

func Test_chain_Prepend(t *testing.T) {
	t1 := tagMiddleware("t1")
	t2 := tagMiddleware("t2")
	t3 := tagMiddleware("t3")

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	chain := New(t3)

	chain.Prepend(t1, t2)
	c := chain.Then(testApp)

	c.ServeHTTP(w, r)

	assert.Equal(t, "t1t2t3app", w.Body.String(),
		"Then does not order handlers correctly")
}

func Test_chain_ThenFunc(t *testing.T) {
	t1 := tagMiddleware("t1")
	t2 := tagMiddleware("t2")
	t3 := tagMiddleware("t3")

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	c := New(t1, t2, t3).ThenFunc(testApp)

	c.ServeHTTP(w, r)

	assert.Equal(t, "t1t2t3app", w.Body.String(),
		"Then does not order handlers correctly")
}
