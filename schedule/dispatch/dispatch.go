package dispatch

import (
	"context"
	"hash/fnv"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/schedule"
)

const (
	DefaultWorkers = 2113
	DefaultJobSize = 64
)

type Dispatcher struct {
	maxWorkers int
	workers    []*Worker
	opt        *Option
	ctx        context.Context
}

func New(ctx context.Context, opts ...OptionFunc) schedule.Scheduler {
	opt := &Option{
		maxWorker: DefaultWorkers,
	}

	for _, fn := range opts {
		fn(opt)
	}

	d := &Dispatcher{
		maxWorkers: opt.maxWorker,
		opt:        opt,
		ctx:        ctx,
	}

	d.workers = make([]*Worker, d.maxWorkers)

	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, DefaultJobSize)
		worker.start(d.ctx)
		d.workers[i] = worker
	}

	return d
}

func (d *Dispatcher) Schedule(j *schedule.Job) {
	key := j.Key

	h := fnv.New32a()
	h.Write([]byte(key))
	sum := h.Sum32()
	slot := int(sum) % d.maxWorkers

	woker := d.workers[slot]
	select {
	case woker.jobQueue <- j:
		return
	default: //工作队列塞满了，丢弃请求
		alog.Fatal("%s|%s|%d|%d|job queue is full, drop it\n", key, j.String(), slot, woker.execs)
		alog.Fatal("%d|work ring %s\n", slot, woker.ring.String())
		return
	}
}
