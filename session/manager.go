package session

import (
	"context"
	"sync"
	"time"

	"github.com/lightmen/nami/pkg/safe"
)

type Manager struct {
	sessions sync.Map
	ttlMilli int64 //session超时删除时间，单位豪秒
}

var gMgr *Manager

func init() {
	gMgr = NewManager()
}

type ManagerOption func(mgr *Manager)

func WithTTL(ttlMilli int64) ManagerOption {
	return func(mgr *Manager) {
		mgr.ttlMilli = ttlMilli
	}
}

func NewManager(opts ...ManagerOption) *Manager {
	mgr := &Manager{
		ttlMilli: 180000, //默认180秒,超过该时间的session会从内存中清掉
	}

	for _, op := range opts {
		op(mgr)
	}

	safe.Go(mgr.loop)
	return mgr
}

func (mgr *Manager) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		_ = t
		mgr.cleanSession()
	}
}

func (mgr *Manager) cleanSession() {
	const (
		maxScanCount = 1024
		maxDelCount  = 256
	)
	delCount := 0
	scanCount := 0

	nowMilli := time.Now().UnixMilli()

	mgr.sessions.Range(func(key, val any) bool {
		s := val.(*session)
		lastTime := s.GetLastTime()
		if nowMilli-lastTime >= mgr.ttlMilli {
			mgr.sessions.Delete(key)
			delCount++
		}

		if delCount >= maxDelCount {
			return false
		}

		scanCount++
		if scanCount >= maxScanCount {
			return false
		}

		return true
	})
}

// Get 根据uid返回玩家的session
func (mgr *Manager) Get(uid string) (s *session, ok bool) {
	val, ok := mgr.sessions.Load(uid)
	if !ok {
		return
	}

	s = val.(*session)
	s.SetLastTime(time.Now().UnixMilli())
	return
}

func (mgr *Manager) Set(s *session) {
	if s == nil {
		return
	}
	mgr.sessions.Store(s.ID(), s)
}

func Get(uid string) (s *session, ok bool) {
	return gMgr.Get(uid)
}

func Set(s *session) {
	gMgr.Set(s)
}

// GetForce 根据uid获取session，如果不存在则创建一个
func GetForce(ctx context.Context, uid string) (s Session) {
	s, ok := Get(uid)
	if ok {
		return s
	}

	if gFactory != nil {
		s = gFactory.Get(ctx, uid)
		return
	}

	s = New(uid)
	return
}
