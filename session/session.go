package session

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Session interface {
	ID() string //返回session ID
	Set(key, val any)
	Get(key any) (val any, ok bool)
}

// Factory 创建Session函数
type Factory interface {
	Get(ctx context.Context, id string) Session
}

var gFactory Factory

func RegisterFactory(factory Factory) {
	gFactory = factory
}

type session struct {
	id       string       // binding user id
	lastTime atomic.Int64 // last heartbeat time
	data     sync.Map
}

// New returns a new session instance
func New(id string) *session {
	s := &session{
		id: id,
	}

	s.SetLastTime(time.Now().UnixMilli())

	Set(s)

	return s
}

func (s *session) Get(key any) (val any, ok bool) {
	return s.data.Load(key)
}

func (s *session) Set(key, val any) {
	s.data.Store(key, val)
}

// UID returns uid that bind to current session
func (s *session) ID() string {
	return s.id
}

func (s *session) GetLastTime() int64 {
	return s.lastTime.Load()
}

func (s *session) SetLastTime(lastTime int64) {
	s.lastTime.Store(lastTime)
}
