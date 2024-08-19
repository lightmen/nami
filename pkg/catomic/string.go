package catomic

import "sync/atomic"

// String is an atomic type-safe wrapper for string values
type String struct {
	_ nocmp // disallow non-atomic comparison
	v atomic.Value
}

var _zeroString string

func NewString(str string) *String {
	s := &String{}
	if str != _zeroString {
		s.Store(str)
	}

	return s
}

func (s *String) Store(str string) {
	s.v.Store(str)
}

func (s *String) Load() string {
	if v := s.v.Load(); v != nil {
		return v.(string)
	}

	return _zeroString
}
