package dispatch

import (
	"strconv"

	"github.com/lightmen/nami/schedule"
)

const (
	maxRingSize = 32
)

type Ring struct {
	Data []*schedule.Job
	idx  int
}

func NewRing(size int) *Ring {
	if size > maxRingSize || size < 1 {
		size = maxRingSize
	}

	return &Ring{
		Data: make([]*schedule.Job, size),
		idx:  -1,
	}
}

func (q *Ring) Push(j *schedule.Job) {
	q.idx++
	if q.idx >= len(q.Data) {
		q.idx = 0
	}
	q.Data[q.idx] = j
}

func (q *Ring) String() string {
	str := "[" + strconv.Itoa(q.idx) + "]:"
	for i := 0; i < len(q.Data); i++ {
		if i == q.idx {
			str += "**"
		}

		str += q.Data[i].String()

		if i == q.idx {
			str += "**"
		}

		if i < len(q.Data)-1 {
			str += ","
		}
	}

	return str
}
