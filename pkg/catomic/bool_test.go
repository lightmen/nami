package catomic

import "testing"

func TestBool(t *testing.T) {
	b := &Bool{}

	if b.Load() {
		t.Fatal("test failed")
	}

	b.Store(true)
	if !b.Load() {
		t.Fatal("test failed")
	}

	b.Store(false)
	if b.Load() {
		t.Fatal("test failed")
	}
}
