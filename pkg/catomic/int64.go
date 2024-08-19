package catomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

// Int64 is an atomic wrapper around int64.
type Int64 struct {
	_ nocmp // disallow non-atomic comparison

	v int64
}

// NewInt64 creates a new Int64.
func NewInt64(val int64) *Int64 {
	return &Int64{v: val}
}

// Load atomically loads the wrapped value.
func (i *Int64) Load() int64 {
	return atomic.LoadInt64(&i.v)
}

// Add atomically adds to the wrapped int64 and returns the new value.
func (i *Int64) Add(delta int64) int64 {
	return atomic.AddInt64(&i.v, delta)
}

// Sub atomically subtracts from the wrapped int64 and returns the new value.
func (i *Int64) Sub(delta int64) int64 {
	return atomic.AddInt64(&i.v, -delta)
}

// Inc atomically increments the wrapped int64 and returns the new value.
func (i *Int64) Inc() int64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int64 and returns the new value.
func (i *Int64) Dec() int64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int64) CAS(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Int64) Store(val int64) {
	atomic.StoreInt64(&i.v, val)
}

// Swap atomically swaps the wrapped int64 and returns the old value.
func (i *Int64) Swap(val int64) (old int64) {
	return atomic.SwapInt64(&i.v, val)
}

// MarshalJSON encodes the wrapped int64 into JSON.
func (i *Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped int64.
func (i *Int64) UnmarshalJSON(b []byte) error {
	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (i *Int64) String() string {
	v := i.Load()
	return strconv.FormatInt(v, 10)
}
