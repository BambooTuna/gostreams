package stream

import (
	"errors"
	"reflect"
	"sync"
)

type Sink struct {
	mtx sync.RWMutex
	in  *In

	fn interface{}
}

func NewSink(fn interface{}) *Sink {
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

	if err := isValidHandler(s.fn); err != nil {
		return err
	}
	rv := reflect.ValueOf(s.fn)

	if err := s.in.run(); err != nil {
		return err
	}

	go func() {
		for arg := range s.in.queue {
			rv.Call(arg)
		}
	}()

	return nil
}
