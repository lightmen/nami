package hash

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testKey = "test_hash"

func makeHash(n int) *Hash {
	hash := New()

	for i := 1; i <= n; i++ {
		hkey := fmt.Sprintf("key_%d", i)
		hval := fmt.Sprintf("val_%d", i)
		hash.HSet(testKey, hkey, hval)
	}

	return hash
}

type kv struct {
	k, v string
}

func TestHash_struct(t *testing.T) {
	hash := New()
	rec := &kv{"a", "b"}
	r1 := hash.HSet(testKey, "c", rec)
	assert.Equal(t, r1, 1)
	r2 := hash.HGet(testKey, "c")
	assert.Equal(t, r2, rec)
}

func TestHash_bytes(t *testing.T) {
	hash := New()

	hash.HSet(testKey, "a", []byte("hash_1"))
	hash.HSet(testKey, "b", []byte("hash_2"))
	hash.HSet(testKey, "c", []byte("hash_3"))

	r1 := hash.HGet(testKey, "a")
	assert.Equal(t, r1, []byte("hash_1"))
	r2 := hash.HGet(testKey, "b")
	assert.Equal(t, r2, []byte("hash_2"))
	r3 := hash.HGet(testKey, "c")
	assert.Equal(t, r3, []byte("hash_3"))

}

func TestHash_HSet(t *testing.T) {
	hash := makeHash(3)
	r1 := hash.HSet(testKey, "d", "123")
	assert.Equal(t, r1, 1)
	r2 := hash.HSet(testKey, "d", "123")
	assert.Equal(t, r2, 0)
	r3 := hash.HSet(testKey, "e", "234")
	assert.Equal(t, r3, 1)
}

func TestHash_HSetNx(t *testing.T) {
	hash := makeHash(3)
	r2 := hash.HSetNx(testKey, "bar", "old")
	assert.Equal(t, r2, 1)
	r3 := hash.HSetNx(testKey, "bar", "old")
	assert.Equal(t, r3, 0)
}

func TestHash_HGet(t *testing.T) {
	hash := makeHash(3)
	val := hash.HGet(testKey, "key_1")
	assert.Equal(t, "val_1", val.(string))
	valNotExist := hash.HGet(testKey, "m")
	assert.Equal(t, nil, valNotExist)
}

func TestHash_HDel(t *testing.T) {
	hash := makeHash(3)

	// delete existed key,return 1
	res := hash.HDel(testKey, "key_1")
	assert.Equal(t, 1, res)

	//delete non existing key,return 0
	res = hash.HDel(testKey, "key_9")
	assert.Equal(t, 0, res)
}

func TestHash_HExists(t *testing.T) {
	hash := makeHash(3)

	// key exist
	exist := hash.HExists(testKey, "key_1")
	assert.Equal(t, true, exist)

	// key does non exist
	not := hash.HExists(testKey, "r")
	assert.Equal(t, false, not)

}

func TestHash_HVals(t *testing.T) {
	hash := makeHash(3)
	values := hash.HVals(testKey)
	for i, v := range values {
		assert.Equal(t, ("val_" + strconv.Itoa(i+1)), v)
	}
}

func TestHash_HLen(t *testing.T) {
	hash := makeHash(3)
	assert.Equal(t, 3, hash.HLen(testKey))
}

func TestHash_Keys(t *testing.T) {
	hash := New()
	hash.HSet("k1", "a", []byte("hash_1"))

	n := hash.Keys()
	assert.Equal(t, 1, len(n))
}
