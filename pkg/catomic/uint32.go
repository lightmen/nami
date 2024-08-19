package catomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

// Uint32 is an atomic wrapper around uint32.
type Uint32 struct {
	_ nocmp // disallow non-atomic comparison

	v uint32
}

// NewUint32 creates a new Uint32.
func NewUint32(val uint32) *Uint32 {
	return &Uint32{v: val}
}

// Load atomically loads the wrapped value.
func (i *Uint32) Load() uint32 {
	return atomic.LoadUint32(&i.v)
}

// Add atomically adds to the wrapped uint32 and returns the new value.
func (i *Uint32) Add(delta uint32) uint32 {
	return atomic.AddUint32(&i.v, delta)
}

// Sub atomically subtracts from the wrapped uint32 and returns the new value.
func (i *Uint32) Sub(delta uint32) uint32 {
	return atomic.AddUint32(&i.v, ^(delta - 1))
}

// Inc atomically increments the wrapped uint32 and returns the new value.
func (i *Uint32) Inc() uint32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint32 and returns the new value.
func (i *Uint32) Dec() uint32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint32) CAS(old, new uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Uint32) Store(val uint32) {
	atomic.StoreUint32(&i.v, val)
}

// Swap atomically swaps the wrapped uint32 and returns the old value.
func (i *Uint32) Swap(val uint32) (old uint32) {
	return atomic.SwapUint32(&i.v, val)
}

// MarshalJSON encodes the wrapped uint32 into JSON.
func (i *Uint32) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}

// UnmarshalJSON decodes JSON into the wrapped uint32.
func (i *Uint32) UnmarshalJSON(b []byte) error {
	var v uint32
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}

// String encodes the wrapped value as a string.
func (i *Uint32) String() string {
	v := i.Load()
	return strconv.FormatUint(uint64(v), 10)
}
