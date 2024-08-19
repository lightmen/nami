package catomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

// Uint64 is an atomic wrapper around uint64.
type Uint64 struct {
	_ nocmp // disallow non-atomic comparison

	v uint64
}

// NewUint64 creates a new Uint64.
func NewUint64(val uint64) *Uint64 {
	return &Uint64{v: val}
}

// Load atomically loads the wrapped value.
func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64(&i.v)
}

// Add atomically adds to the wrapped uint64 and returns the new value.
func (i *Uint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&i.v, delta)
}

// Sub atomically subtracts from the wrapped uint64 and returns the new value.
func (i *Uint64) Sub(delta uint64) uint64 {
	return atomic.AddUint64(&i.v, ^(delta - 1))
}

// Inc atomically increments the wrapped uint64 and returns the new value.
func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint64 and returns the new value.
func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint64) CAS(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Uint64) Store(val uint64) {
	atomic.StoreUint64(&i.v, val)
}

// Swap atomically swaps the wrapped uint64 and returns the old value.
func (i *Uint64) Swap(val uint64) (old uint64) {
	return atomic.SwapUint64(&i.v, val)
}

// MarshalJSON encodes the wrapped uint64 into JSON.
func (i *Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped uint64.
func (i *Uint64) UnmarshalJSON(b []byte) error {
	var v uint64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (i *Uint64) String() string {
	v := i.Load()
	return strconv.FormatUint(uint64(v), 10)
}
