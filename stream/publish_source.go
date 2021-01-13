package stream

import (
	"fmt"
	"reflect"
	"sync"
)

type sourceImpl struct {
	mtx      sync.RWMutex
	outQueue chan []reflect.Value
}

func NewSource() *sourceImpl {
	return &sourceImpl{
		outQueue: make(chan []reflect.Value, 30),
	}
}

func (s *sourceImpl) Publish(arg interface{}) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	select {
	case s.outQueue <- buildHandlerArg(arg):
	default:
		fmt.Println("Overflowing: The element disappears", arg)
	}
}

func (s *sourceImpl) Via(flow Flow) {
	flow.attachInBound(&In{
		queue: s.outQueue,
		run: func() error {
			return nil
		},
	})
}

func (s *sourceImpl) To(sink Sink) {
	sink.attachInBound(&In{
		queue: s.outQueue,
		run: func() error {
			return nil
		},
	})
}
