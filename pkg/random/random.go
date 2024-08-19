package random

import (
	"math/rand"
	"unsafe"
)

const btstable = "0123456789abcdefghjmnqrstABCDEFGHJKLMNPQRSTUVWXYZ"

func Bytes(size int) []byte {
	if size > 128 {
		size = 128
	}

	var bts [128]byte
	ln := len(btstable)

	for i := 0; i < size; i++ {
		bts[i] = btstable[rand.Intn(ln)]
	}

	return bts[:size]
}

func String(size int) string {
	bts := Bytes(size)

	return unsafe.String(&bts[0], size)
}

func RandomInRange(min, max int64) int64 {
	if min >= max {
		return min
	}

	return rand.Int63n(max-min) + min
}

func RandomSubset[T int | int32 | int64](min, max T, n int) []T {
	if min >= max {
		return []T{}
	}

	if n <= 0 {
		return []T{}
	}

	if int64(n) > int64(max)-int64(min) {
		return []T{}
	}
	randList := rand.Perm(int(max) - int(min))
	ret := make([]T, n)
	for i := 0; i < n; i++ {
		ret[i] = T(randList[i]) + min
	}
	return ret
}

// IsRanded 给定权重，判断是否能随机到，默认是万分比随机
func IsRanded(weight int32, opts ...Option) bool {
	if weight <= 0 {
		return false
	}

	opt := &option{
		randVal: 10000,
	}
	for _, o := range opts {
		o(opt)
	}

	if opt.randVal <= 0 {
		return false
	}

	val := rand.Int31n(opt.randVal)

	return val < weight
}

func RandomElement[T any](arr []T) T {

	// 随机生成一个索引
	index := rand.Intn(len(arr))

	// 返回该索引处的元素
	return arr[index]
}

// Array 随机打乱arr的元素
func Array[T any](arr []T) []T {
	if len(arr) == 0 {
		return arr
	}

	idxs := rand.Perm(len(arr))

	for idx, val := range idxs {
		arr[idx], arr[val] = arr[val], arr[idx]
	}

	return arr
}
