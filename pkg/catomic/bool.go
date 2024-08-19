package catomic

import "sync/atomic"

type Bool struct {
	_ nocmp
	v uint32
}

var _zeroBool bool

func NewBool(val bool) *Bool {
	b := &Bool{}
	if val != _zeroBool {
		b.Store(val)
	}

	return b
}

func (b *Bool) Store(val bool) {
	atomic.StoreUint32(&b.v, bool2Uint(val))
}

func (b *Bool) Load() bool {
	val := atomic.LoadUint32(&b.v)
	return uint2Bool(val)
}

func (b *Bool) Swap(val bool) bool {
	old := atomic.SwapUint32(&b.v, bool2Uint(val))
	return uint2Bool(old)
}

func (b *Bool) CAS(old, new bool) (swapped bool) {
	return atomic.CompareAndSwapUint32(&b.v, bool2Uint(old), bool2Uint(new))
}
