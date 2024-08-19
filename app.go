package nami

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport"
	"github.com/rs/xid"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts     *options
	ctx      context.Context
	cancel   context.CancelFunc
	instance *registry.Instance
	lk       sync.RWMutex
}

var gApp *App

func GetInfo() AppInfo {
	return GetApp()
}

func GetApp() *App {
	return gApp
}

func GetContext() context.Context {
	return GetApp().Context()
}

func New(opts ...Option) (a *App, err error) {
	id := xid.New()
	o := &options{
		ctx:  context.Background(),
		id:   id.String(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	for _, opt := range opts {
		opt(o)
	}

	ctx, cancel := context.WithCancel(o.ctx)

	a = &App{
		cancel: cancel,
		opts:   o,
	}

	a.ctx = NewContext(ctx, a)

	if gApp == nil {
		gApp = a
	}

	return
}

func (a *App) Run() (err error) {
	alog.Info("app %s:%s start", a.opts.name, a.opts.id)

	instance, err := a.buildInstance()
	if err != nil {
		return
	}
	a.updateInstance(instance)

	sctx := a.ctx

	group, ctx := errgroup.WithContext(sctx)

	for _, fn := range a.opts.beforeStart {
		if err = fn(sctx); err != nil {
			return
		}
	}

	err = a.startServer(group, ctx)
	if err != nil {
		return
	}

	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
		defer rcancel()

		if err = a.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}

	for _, fn := range a.opts.afterStart {
		if err = fn(sctx); err != nil {
			return
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	group.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				alog.InfoCtx(ctx, "ctx done, app exit")
				return nil
			case sig := <-c:
				if a.opts.sigFunc != nil && !a.opts.sigFunc(sig) {
					alog.InfoCtx(ctx, "got signal %s, but can not exit", sig.String())
					continue
				}
				alog.InfoCtx(ctx, "got signal %s exit", sig.String())
				return a.Stop()
			}
		}
	})

	if err = group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		alog.InfoCtx(ctx, "group stopped, err: %s", err.Error())
		return
	}

	for _, fn := range a.opts.afterStop {
		err = fn(sctx)
	}

	alog.InfoCtx(ctx, "app %s:%s end", a.opts.name, a.opts.id)
	return err
}

func (a *App) startServer(group *errgroup.Group, ctx context.Context) (err error) {
	wg := sync.WaitGroup{}

	for _, s := range a.opts.servers {
		srv := s

		group.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(NewContext(a.opts.ctx, a), 10*time.Second)
			defer cancel()
			return srv.Stop(stopCtx)
		})

		wg.Add(1)
		group.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}

	wg.Wait()
	return
}

func (a *App) buildInstance() (*registry.Instance, error) {
	endpoints := make([]string, 0)
	for _, srv := range a.opts.servers {
		if r, ok := srv.(transport.Endpointer); ok {
			e, err := r.Endpoint()
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, e.String())
		}
	}

	return &registry.Instance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		MetaData:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

func (a *App) updateInstance(instance *registry.Instance) {
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()
}

func (a *App) getInstance() *registry.Instance {
	a.lk.RLock()
	instance := a.instance
	a.lk.RUnlock()

	return instance
}

func (a *App) Stop() (err error) {
	alog.Info("app %s:%s stop", a.opts.name, a.opts.id)

	instance := a.getInstance()

	if a.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(a.ctx, a), 5*time.Second)
		defer cancel()
		if err = a.opts.registrar.Unregister(ctx, instance); err != nil {
			alog.InfoCtx(ctx, "app %s:%s Unregister error: %s", a.opts.name, a.opts.id, err.Error())
			return
		}
	}

	if a.cancel != nil {
		a.cancel()
	}

	return
}

func (a *App) Context() context.Context {
	return a.ctx
}
