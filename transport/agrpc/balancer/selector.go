package balancer

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type Matcher interface {
	Match(ctx context.Context) bool
}

type Picker interface {
	Matcher
	base.PickerBuilder
	balancer.Picker
	Name() string
}

type selector struct {
	sync.RWMutex
	pickers map[string]Picker
}

func (s *selector) Register(picker Picker) {
	name := picker.Name()

	s.Lock()
	defer s.Unlock()

	if _, ok := s.pickers[name]; ok {
		panic(fmt.Sprintf("duplicate register selector %s", name))
	}

	s.pickers[name] = picker
}

func (s *selector) Get(ctx context.Context) Picker {
	s.RLock()
	defer s.RUnlock()

	for _, picker := range s.pickers {
		if picker.Match(ctx) {
			return picker
		}
	}

	return nil
}

func (s *selector) Update(info base.PickerBuildInfo) {
	s.RLock()
	defer s.RUnlock()

	for _, builder := range s.pickers {
		builder.Build(info)
	}
}
