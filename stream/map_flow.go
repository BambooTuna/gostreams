package stream

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type mapFlow struct {
	mtx sync.RWMutex

	in *In

	outQueue chan []reflect.Value

	fn interface{}
}

func NewMapFlow(fn interface{}) *mapFlow {
	return &mapFlow{
		outQueue: make(chan []reflect.Value, 30),
		fn:       fn,
	}
}

func (f *mapFlow) attachInBound(in *In) {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	f.in = in
}

func (f *mapFlow) run() error {
	if f.in == nil {
		return errors.New("please attach source ")
	}

	if err := isValidHandler(f.fn); err != nil {
		return err
	}
	rv := reflect.ValueOf(f.fn)

	if err := f.in.run(); err != nil {
		return err
	}

	go func() {
		for arg := range f.in.queue {
			select {
			case f.outQueue <- rv.Call(arg):
			default:
				fmt.Println("Overflowing: The element disappears.", arg)
			}
		}
	}()

	return nil
}

func (f *mapFlow) Via(flow Flow) {
	flow.attachInBound(&In{
		queue: f.outQueue,
		run:   f.run,
	})
}

func (f *mapFlow) To(sink Sink) {
	sink.attachInBound(&In{
		queue: f.outQueue,
		run:   f.run,
	})
}
