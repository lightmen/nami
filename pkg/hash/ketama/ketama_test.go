package ketama

import "testing"

func TestKetama(t *testing.T) {
	k := New(Replicas(9))

	nodes := []string{"node1", "node2", "node3", "node4"}

	for _, n := range nodes {
		k.Add(n)
	}

	testKey := "testKey1"
	expect, ok := k.Get(testKey)
	if !ok {
		t.Fatalf("get failed")
	}

	got, _ := k.Get(testKey)

	if got != expect {
		t.Fatalf("expect %s, got: %s", expect, got)
	}

	for _, n := range nodes {
		k.Remove(n)
	}

	if !k.IsEmpty() {
		t.Fatalf("is not empty")
	}

	got, ok = k.Get(testKey)
	if ok {
		t.Fatalf("expect false, got true")
	}
}
