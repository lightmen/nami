package nami

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/lightmen/nami/core/log"
	"github.com/lightmen/nami/registry"
	"github.com/lightmen/nami/transport"
)

type App struct {
	opts     *options
	logger   log.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	instance *registry.Instance
	lk       sync.RWMutex
}

func New(opts ...Option) (s *App, err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}

	o := &options{
		id:   id.String(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	for _, opt := range opts {
		opt(o)
	}

	s = &App{
		ctx:    context.Background(),
		opts:   o,
		logger: o.logger,
	}

	return
}

func (a *App) Run() (err error) {
	instance, err := a.buildInstance()
	if err != nil {
		return
	}

	var ctx context.Context
	ctx, a.cancel = context.WithCancel(a.ctx)

	wg := sync.WaitGroup{}

	for _, s := range a.opts.servers {
		srv := s

		go func() {
			<-ctx.Done() //wait for stop signal
			srv.Stop()
		}()

		wg.Add(1)
		go func() {
			wg.Done()

			err = srv.Start()
			if err != nil {
				a.logger.Error("srv.Start error: %s", err.Error())
				return
			}
		}()
	}

	wg.Wait()

	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(a.opts.ctx, 10*time.Second)
		defer rcancel()

		if err = a.opts.registrar.Register(rctx, instance); err != nil {
			return
		}

		a.updateInstance(instance)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)

	select {
	case <-c:
		a.Stop()
	}

	return
}

func (a *App) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
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
