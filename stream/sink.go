package stream

import (
	"errors"
	"sync"
)

type Sink struct {
	mtx sync.RWMutex
	in  *In

	fn func(Item)
}

func NewSink(fn func(Item)) *Sink {
	return &Sink{fn: fn}
}

func (s *Sink) attachInBound(in *In) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	s.in = in
}

func (s *Sink) Run() error {
	if s.in == nil {
		return errors.New("please attach source ")
	}

	if err := s.in.run(); err != nil {
		return err
	}

	go func() {
		for arg := range s.in.queue {
			s.fn(arg)
		}
	}()

	return nil
}
