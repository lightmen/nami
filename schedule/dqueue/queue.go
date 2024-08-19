package dqueue

// 动态队列，worker 数量会自动增长/缩减

import (
	"container/list"
	"context"
	"sync"
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/schedule"
)

const (
	// 最大空闲 worker 数量
	maxIdleWorkers = 128

	// 每次回收 worker 数量
	workerRecycleNum = 64

	// 每次回收 jobUnion 数量
	jobUnionRecycleNum = 128
)

var _ schedule.Scheduler = (*dQueue)(nil)

type dQueue struct {
	lock       sync.Mutex
	opt        *option
	maxWorkers int // worker 数量超过上限会告警，但是仍然可以增长
	workingNum int // 当前在工作中的 worker 数量
	ctx        context.Context
	ch         unionChan

	waitWorkers *list.List           // 空闲等待的 worker 队列
	unMap       map[string]*jobUnion // 用来快速查找 jobUnion 的 map
}

func New(ctx context.Context, opts ...OptionFunc) schedule.Scheduler {
	opt := &option{
		maxWorker: maxIdleWorkers,
	}

	for _, fn := range opts {
		fn(opt)
	}

	q := &dQueue{
		opt:         opt,
		maxWorkers:  opt.maxWorker,
		ctx:         ctx,
		ch:          make(unionChan),
		waitWorkers: list.New(),
		unMap:       make(map[string]*jobUnion, 5000),
	}

	for i := 0; i < maxIdleWorkers; i++ {
		w := newWorker(ctx, q)
		q.waitWorkers.PushBack(w)
	}

	go q.wait()

	return q
}

func (q *dQueue) wait() {
	tk := time.NewTicker(time.Second * 30)
	defer tk.Stop()

	for {
		select {
		case un := <-q.ch:
			w := q.getWorker(un)
			w.ch <- un

		case <-q.ctx.Done():
			return

		case <-tk.C:
			q.free()
		}
	}

}

func (q *dQueue) getWorker(un *jobUnion) *worker {
	q.lock.Lock()
	defer q.lock.Unlock()

	// jobUnion 有 worker 正在执行它的 job
	if un.w != nil {
		return un.w
	}

	ele := q.waitWorkers.Front()
	if ele == nil {
		q.grow()
		ele = q.waitWorkers.Front()
	}
	q.waitWorkers.Remove(ele)
	q.workingNum++

	w := ele.Value.(*worker)
	un.w = w

	return un.w
}

func (q *dQueue) putWorker(w *worker, un *jobUnion) (idle bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if un.Len() > 0 {
		return false
	}

	q.waitWorkers.PushFront(w)
	q.workingNum--
	un.w = nil
	return true
}

// 动态增加 worker;
// 调用者必须持有 q.lock.Lock()
func (q *dQueue) grow() {
	for i := 0; i < maxIdleWorkers; i++ {
		w := newWorker(q.ctx, q)
		q.waitWorkers.PushBack(w)
	}

	wn := q.waitWorkers.Len()
	ln := q.workingNum + wn
	if ln > q.maxWorkers {
		alog.Fatal("worker queue %d exceeded maxWorkers: %d, waiting workers: %d", ln, q.maxWorkers, wn)
	}
}

// 回收多余空转的 worker
func (q *dQueue) free() {
	q.lock.Lock()
	defer q.lock.Unlock()

	// 等待队列保留 maxIdleWorkers 数量的空闲 worker
	// 每次回收动作执行上限为 workerRecycleNum
	num := q.waitWorkers.Len() - maxIdleWorkers
	if num > 0 {
		if num > workerRecycleNum {
			num = workerRecycleNum
		}

		for i := 0; i < num; i++ {
			ele := q.waitWorkers.Back()
			if ele == nil {
				break
			}

			q.waitWorkers.Remove(ele)
			w := ele.Value.(*worker)
			w.stop()
		}
	}

	var cnt int
	for k, un := range q.unMap {
		if un.w == nil && un.Len() == 0 {
			delete(q.unMap, k)
		}
		cnt++

		if cnt > jobUnionRecycleNum {
			break
		}
	}
}

func (q *dQueue) Schedule(j *schedule.Job) {
	q.lock.Lock()

	un, ok := q.unMap[j.Key]
	if !ok {
		un = newJobUnion(j.Key)
		q.unMap[j.Key] = un
	}
	empty := un.PushBack(j)

	q.lock.Unlock()

	// jobs 从空转为非空状态，触发执行
	if empty {
		q.ch <- un
	}
}
