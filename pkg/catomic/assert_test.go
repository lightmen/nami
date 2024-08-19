package catomic

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Marks the test as failed if the error cannot be cast into the provided type
// with errors.As.
//
//	assertErrorAsType(t, err, new(ErrFoo))
func assertErrorAsType(t *testing.T, err error, typ interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	return assert.True(t, errors.As(err, typ), msgAndArgs...)
}

func assertErrorJSONUnmarshalType(t *testing.T, err error, msgAndArgs ...interface{}) bool {
	t.Helper()

	return assertErrorAsType(t, err, new(*json.UnmarshalTypeError), msgAndArgs...)
}
