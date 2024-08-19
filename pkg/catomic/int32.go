package catomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

// Int32 is an atomic wrapper around Int32.
type Int32 struct {
	_ nocmp // disallow non-atomic comparison

	v int32
}

// NewInt32 creates a new Int32.
func NewInt32(val int32) *Int32 {
	return &Int32{v: val}
}

// Load atomically loads the wrapped value.
func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}

// Add atomically adds to the wrapped Int32 and returns the new value.
func (i *Int32) Add(delta int32) int32 {
	return atomic.AddInt32(&i.v, delta)
}

// Sub atomically subtracts from the wrapped Int32 and returns the new value.
func (i *Int32) Sub(delta int32) int32 {
	return atomic.AddInt32(&i.v, -delta)
}

// Inc atomically increments the wrapped Int32 and returns the new value.
func (i *Int32) Inc() int32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped Int32 and returns the new value.
func (i *Int32) Dec() int32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int32) CAS(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Int32) Store(val int32) {
	atomic.StoreInt32(&i.v, val)
}

// Swap atomically swaps the wrapped Int32 and returns the old value.
func (i *Int32) Swap(val int32) (old int32) {
	return atomic.SwapInt32(&i.v, val)
}

// MarshalJSON encodes the wrapped Int32 into JSON.
func (i *Int32) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped Int32.
func (i *Int32) UnmarshalJSON(b []byte) error {
	var v int32
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (i *Int32) String() string {
	v := i.Load()
	return strconv.FormatInt(int64(v), 10)
}
