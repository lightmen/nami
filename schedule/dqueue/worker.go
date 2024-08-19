package dqueue

import (
	"context"

	"github.com/lightmen/nami/pkg/safe"
)

type worker struct {
	ch     unionChan
	ctx    context.Context
	cancel context.CancelFunc
	q      *dQueue
}

func newWorker(ctx context.Context, q *dQueue) *worker {
	// ch 缓存给 1 是为了在极端情况下，run() 处理了最后一个 job 的同时，插入了一个新的 job
	// 此时 dQueue.wait() 会阻塞在
	// case e := <-q.ch:
	//     w := q.getWorker(e)
	//     w.ch <- e
	// 给 1 个缓存可以防止这种情况下阻塞
	w := &worker{
		ch: make(unionChan, 1),
		q:  q,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)

	w.run()

	return w
}

func (w *worker) run() {
	go func() {
		for {
			select {
			case un := <-w.ch:
				idle := false
				for !idle {
					j, last := un.PopFront()
					for j != nil {
						safe.Func(func() {
							j.Jobber(j)
						})

						if last {
							break
						}
						j, last = un.PopFront()
					}

					idle = w.q.putWorker(w, un)
				}

			case <-w.ctx.Done():
				w.ch = nil
				w.q = nil
				return
			}
		}
	}()
}

func (w *worker) stop() {
	w.cancel()
}
