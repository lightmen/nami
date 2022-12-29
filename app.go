package nami

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/lightmen/nami/log"
)

type App struct {
	opts   *options
	logger log.Logger
	ctx    context.Context
	cancel context.CancelFunc
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
