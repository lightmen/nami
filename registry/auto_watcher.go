package registry

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/lightmen/nami/alog"
)

type AutoWatcher struct {
	watchName string
	dis       Discovery
	timeout   time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
	w         Watcher
	sync.RWMutex
	fnList []func([]*Instance)
}

func NewAutoWatcher(watchName string, dis Discovery, fn func([]*Instance)) (*AutoWatcher, error) {
	aw := &AutoWatcher{
		watchName: watchName,
		dis:       dis,
		timeout:   time.Second * 3,
	}

	aw.Register(fn)
	if err := aw.work(); err != nil {
		return nil, err
	}

	return aw, nil
}

func (aw *AutoWatcher) Register(fn func([]*Instance)) {
	aw.Lock()
	aw.fnList = append(aw.fnList, fn)
	aw.Unlock()

}

func (aw *AutoWatcher) work() (err error) {
	watchRes := &struct {
		err error
		w   Watcher
	}{}

	done := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		w, err := aw.dis.Watch(ctx, aw.watchName)
		watchRes.w = w
		watchRes.err = err
		close(done)
	}()

	select {
	case <-done:
		err = watchRes.err
	case <-time.After(aw.timeout):
		err = errors.New("discovery create watcher overtime")
	}
	if err != nil {
		cancel()
		return
	}

	aw.w = watchRes.w
	aw.ctx = ctx
	aw.cancel = cancel

	go aw.watch()

	return
}

func (aw *AutoWatcher) watch() {
	for {
		select {
		case <-aw.ctx.Done():
			return
		default:
		}
		ins, err := aw.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			alog.Error("[resolver][%s] Failed to watch discovery endpoint: %v", aw.watchName, err)
			time.Sleep(time.Second)
			continue
		}
		aw.onChange(ins)
	}
}

func (aw *AutoWatcher) onChange(insList []*Instance) {
	aw.RLock()
	defer aw.RUnlock()

	for _, fn := range aw.fnList {
		// safe.Go(func() {
		// 	fn(ins)
		// })
		fn(insList)
	}
}

func (aw *AutoWatcher) Close() {
	aw.cancel()
	err := aw.w.Stop()
	if err != nil {
		alog.Error("[resolver] failed to watch top: %s", err)
	}
}
