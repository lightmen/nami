package ketama

import (
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

type HashFunc func(data []byte) uint32

const (
	DefaultReplicas = 2999
	Salt            = "n*@if09^Ig3h"
)

func DefaultHash(data []byte) uint32 {
	f := fnv.New32()
	f.Write(data)
	return f.Sum32()
}

type Ketama struct {
	sync.RWMutex
	hash     HashFunc
	replicas int
	keys     []int //  Sorted keys
	hashMap  map[int]string
}

func New(opts ...Option) *Ketama {
	k := &Ketama{
		replicas: DefaultReplicas,
		hash:     DefaultHash,
		hashMap:  make(map[int]string),
	}

	for _, op := range opts {
		op(k)
	}

	return k
}

func (h *Ketama) IsEmpty() bool {
	h.Lock()
	defer h.Unlock()

	return len(h.keys) == 0
}

func (h *Ketama) Add(nodes ...string) {
	h.Lock()
	defer h.Unlock()

	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			key := int(h.hash([]byte(Salt + strconv.Itoa(i) + node)))

			if _, ok := h.hashMap[key]; !ok {
				h.keys = append(h.keys, key)
			}
			h.hashMap[key] = node
		}
	}
	sort.Ints(h.keys)
}

func (h *Ketama) Remove(nodes ...string) {
	h.Lock()
	defer h.Unlock()

	deletedKey := make([]int, 0)
	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			key := int(h.hash([]byte(Salt + strconv.Itoa(i) + node)))

			if _, ok := h.hashMap[key]; ok {
				deletedKey = append(deletedKey, key)
				delete(h.hashMap, key)
			}
		}
	}
	if len(deletedKey) > 0 {
		h.deleteKeys(deletedKey)
	}
}

func (h *Ketama) deleteKeys(deletedKeys []int) {
	sort.Ints(deletedKeys)

	index := 0
	count := 0
	for _, key := range deletedKeys {
		for ; index < len(h.keys); index++ {
			h.keys[index-count] = h.keys[index]

			if key == h.keys[index] {
				count++
				index++
				break
			}
		}
	}

	for ; index < len(h.keys); index++ {
		h.keys[index-count] = h.keys[index]
	}

	h.keys = h.keys[:len(h.keys)-count]
}

func (h *Ketama) Get(key string) (string, bool) {
	if h.IsEmpty() {
		return "", false
	}

	hash := int(h.hash([]byte(key)))

	h.RLock()
	defer h.RUnlock()

	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})

	if idx == len(h.keys) {
		idx = 0
	}
	str, ok := h.hashMap[h.keys[idx]]
	return str, ok
}
