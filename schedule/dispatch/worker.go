package dispatch

import (
	"context"

	"github.com/lightmen/nami/pkg/safe"
	"github.com/lightmen/nami/schedule"
)

type Worker struct {
	id       int
	jobQueue chan *schedule.Job
	execs    int64 //执行次数
	ring     *Ring
}

func NewWorker(id int, jobSize int) *Worker {
	return &Worker{
		id:       id,
		jobQueue: make(chan *schedule.Job, jobSize),
		ring:     NewRing(jobSize),
	}
}

func (w *Worker) start(ctx context.Context) {
	safe.Go(func() {
		for {
			select {
			case job := <-w.jobQueue:
				w.ring.Push(job)
				w.execs++
				fn := func() {
					job.Jobber(job)
				}
				safe.Func(fn)

			case <-ctx.Done():
				return
			}
		}
	})
}
