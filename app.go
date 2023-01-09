package nami

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/lightmen/nami/core/log"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport"
	"golang.org/x/sync/errgroup"
)

type App struct {
	opts     *options
	logger   log.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	instance *registry.Instance
	lk       sync.RWMutex
}

func New(opts ...Option) (a *App, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	o := &options{
		ctx:  context.Background(),
		id:   id.String(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	for _, opt := range opts {
		opt(o)
	}

	a = &App{
		ctx:    context.Background(),
		opts:   o,
		logger: o.logger,
	}

	a.ctx, a.cancel = context.WithCancel(o.ctx)

	return
}

func (a *App) Run() (err error) {
	group, ctx := errgroup.WithContext(a.ctx)

	err = a.startServer(group, ctx)
	if err != nil {
		return
	}

	if err = a.registerInstance(); err != nil {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	group.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})

	if err = group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return
	}

	return nil
}

func (a *App) startServer(group *errgroup.Group, ctx context.Context) (err error) {
	wg := sync.WaitGroup{}

	for _, s := range a.opts.servers {
		srv := s

		group.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
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

func (a *App) registerInstance() (err error) {
	instance, err := a.buildInstance()
	if err != nil {
		return
	}

	a.updateInstance(instance)

	if a.opts.registrar == nil {
		return
	}

	rctx, rcancel := context.WithTimeout(a.opts.ctx, 10*time.Second)
	defer rcancel()

	if err = a.opts.registrar.Register(rctx, instance); err != nil {
		return
	}

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

func (a *App) Stop() (err error) {
	if a.cancel != nil {
		a.cancel()
	}

	return
}
