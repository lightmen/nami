package dqueue

import (
	"container/list"
	"sync"
	"time"

	"github.com/lightmen/nami/schedule"
)

type unionChan chan *jobUnion

type jobUnion struct {
	lock        sync.Mutex
	key         string
	jobs        *list.List
	w           *worker
	lastJobTime time.Time // 最后一个 job 被 worker 取走时间
}

func newJobUnion(key string) *jobUnion {
	return &jobUnion{
		key:         key,
		jobs:        list.New(),
		lastJobTime: time.Now(),
	}
}

// push back 之前，jobs 是否为空
func (u *jobUnion) PushBack(j *schedule.Job) (isEmpty bool) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.jobs.PushBack(j)
	return u.jobs.Len() == 1 && u.w == nil
}

// @last 最后一个 job
func (u *jobUnion) PopFront() (j *schedule.Job, last bool) {
	u.lock.Lock()
	defer u.lock.Unlock()

	ele := u.jobs.Front()
	if ele == nil {
		return
	}
	val := u.jobs.Remove(ele)

	if u.jobs.Len() == 0 {
		last = true
	}

	j = val.(*schedule.Job)
	u.lastJobTime = time.Now()
	return
}

func (u *jobUnion) Len() int {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.jobs.Len()
}
