package pool

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/lightmen/nami/alog"
)

// ErrClosed is the error resulting if the pool is closed via pool.Close().
var ErrClosed = errors.New("pool is closed")

type Pool interface {
	// Get returns a new connection from the pool. Closing the connections puts
	// it back to the Pool. Closing it when the pool is destroyed or full will
	// be counted as an error. we guarantee the conn.Value() isn't nil when conn isn't nil.
	Get() (Conn, error)

	// Close closes the pool and all its connections. After Close() the pool is
	// no longer usable. You can't make concurrent calls Close and Get method.
	// It will be cause panic.
	Close() error

	// Status returns the current status of the pool.
	Status() string
}

type pool struct {
	// atomic, used to get connection random
	index atomic.Uint32

	// atomic, the current physical connection of pool
	current atomic.Int32

	// atomic, the using logic connection of pool
	// logic connection = physical connection * MaxConcurrentStreams
	ref atomic.Int32

	// pool options
	opt options

	// all of created physical connections
	conns []*conn

	// the server address is to create connection.
	address string

	// closed set true when Close is called.
	closed atomic.Int32

	// control the atomic var current's concurrent read write.
	sync.RWMutex
}

// New return a connection pool.
func New(address string, opt ...Option) (Pool, error) {
	if address == "" {
		return nil, errors.New("invalid address settings")
	}

	option := DefaultOptions
	for _, o := range opt {
		o(&option)
	}

	p := &pool{
		opt:     option,
		conns:   make([]*conn, option.MaxActive),
		address: address,
	}
	p.current.Store(int32(option.MaxIdle))

	for i := 0; i < p.opt.MaxIdle; i++ {
		c, err := p.opt.Dial(address)
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("dial is not able to fill the pool: %s", err)
		}
		p.conns[i] = p.wrapConn(c, false)
	}
	alog.Info("new pool success: %v\n", p.Status())

	return p, nil
}

func (p *pool) incrRef() int32 {
	newRef := p.ref.Add(1)
	if newRef == math.MaxInt32 {
		panic(fmt.Sprintf("overflow ref: %d", newRef))
	}
	return newRef
}

func (p *pool) decrRef() {
	newRef := p.ref.Add(-1)
	if newRef < 0 && p.closed.Load() == 0 {
		panic(fmt.Sprintf("negative ref: %d", newRef))
	}
	if newRef == 0 && p.current.Load() > int32(p.opt.MaxIdle) {
		p.Lock()
		if p.ref.Load() == 0 {
			alog.Info("shrink pool: %d ---> %d, decrement: %d, maxActive: %d\n",
				p.current.Load(), p.opt.MaxIdle, p.current.Load()-int32(p.opt.MaxIdle), p.opt.MaxActive)
			p.current.Store(int32(p.opt.MaxIdle))
			p.deleteFrom(p.opt.MaxIdle)
		}
		p.Unlock()
	}
}

func (p *pool) reset(index int) {
	conn := p.conns[index]
	if conn == nil {
		return
	}
	conn.reset()
	p.conns[index] = nil
}

func (p *pool) deleteFrom(begin int) {
	for i := begin; i < p.opt.MaxActive; i++ {
		p.reset(i)
	}
}

// Get see Pool interface.
func (p *pool) Get() (Conn, error) {
	// the first selected from the created connections
	nextRef := p.incrRef()
	p.RLock()
	current := p.current.Load()
	p.RUnlock()
	if current == 0 {
		return nil, ErrClosed
	}
	if nextRef <= current*int32(p.opt.MaxConcurrentStreams) {
		next := p.index.Add(1) % uint32(current)
		return p.conns[next], nil
	}

	// the number connection of pool is reach to max active
	if current == int32(p.opt.MaxActive) {
		// the second if reuse is true, select from pool's connections
		if p.opt.Reuse {
			next := p.index.Add(1) % uint32(current)
			return p.conns[next], nil
		}
		// the third create one-time connection
		c, err := p.opt.Dial(p.address)
		return p.wrapConn(c, true), err
	}

	// the fourth create new connections given back to pool
	p.Lock()
	current = p.current.Load()
	if current < int32(p.opt.MaxActive) && nextRef > current*int32(p.opt.MaxConcurrentStreams) {
		// 2 times the incremental or the remain incremental
		increment := current
		if current+increment > int32(p.opt.MaxActive) {
			increment = int32(p.opt.MaxActive) - current
		}
		var i int32
		var err error
		for i = 0; i < increment; i++ {
			c, er := p.opt.Dial(p.address)
			if er != nil {
				err = er
				break
			}
			p.reset(int(current + i))
			p.conns[current+i] = p.wrapConn(c, false)
		}
		current += i
		alog.Info("grow pool: %d ---> %d, increment: %d, maxActive: %d\n",
			p.current.Load(), current, increment, p.opt.MaxActive)
		p.current.Store(current)
		if err != nil {
			p.Unlock()
			return nil, err
		}
	}
	p.Unlock()
	next := p.index.Add(1) % uint32(current)
	return p.conns[next], nil
}

// Close see Pool interface.
func (p *pool) Close() error {
	p.closed.Store(0)
	p.index.Store(0)
	p.current.Store(0)
	p.ref.Store(0)
	p.deleteFrom(0)
	alog.Info("close pool success: %v\n", p.Status())
	return nil
}

// Status see Pool interface.
func (p *pool) Status() string {
	return fmt.Sprintf("address:%s, index:%d, current:%d, ref:%d. option:%+v",
		p.address, p.index.Load(), p.current.Load(), p.ref.Load(), p.opt)
}
