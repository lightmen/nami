package random

import (
	"math/rand"
	"sync"
)

var _ rand.Source = (*lockedSource)(nil)

type lockedSource struct {
	source rand.Source
	lk     sync.Mutex
}

func NewLockedSource(seed int64) rand.Source {
	ls := &lockedSource{
		source: rand.NewSource(seed),
	}

	return ls
}

func (l *lockedSource) Int63() int64 {
	l.lk.Lock()
	n := l.source.Int63()
	l.lk.Unlock()

	return n
}

func (l *lockedSource) Seed(seed int64) {
	l.lk.Lock()
	l.source.Seed(seed)
	l.lk.Unlock()
}
