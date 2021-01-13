package stream

import (
	"fmt"
	"reflect"
	"sync"
)

type Source struct {
	mtx      sync.RWMutex
	outQueue chan []reflect.Value
}

func NewSource() *Source {
	return &Source{
		outQueue: make(chan []reflect.Value, 30),
	}
}

func (s *Source) Publish(arg interface{}) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	select {
	case s.outQueue <- buildHandlerArg(arg):
	default:
		fmt.Println("Overflowing: The element disappears", arg)
	}
}

func (s *Source) Via(flow *Flow) {
	flow.attachInBound(&In{
		queue: s.outQueue,
		run: func() error {
			return nil
		},
	})
}

func (s *Source) To(sink *Sink) {
	sink.attachInBound(&In{
		queue: s.outQueue,
		run: func() error {
			return nil
		},
	})
}
