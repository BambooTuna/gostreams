package stream

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type groupFlow struct {
	mtx sync.RWMutex

	in *In

	outQueue chan []reflect.Value

	group      int
	groupQueue []interface{}
}

func NewGroupFlow(group int) *groupFlow {
	return &groupFlow{
		outQueue:   make(chan []reflect.Value, 30),
		group:      group,
		groupQueue: make([]interface{}, 0),
	}
}

func (f *groupFlow) attachInBound(in *In) {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	f.in = in
}

func (f *groupFlow) run() error {
	if f.in == nil {
		return errors.New("please attach source ")
	}

	if err := f.in.run(); err != nil {
		return err
	}

	go func() {
		for arg := range f.in.queue {
			f.mtx.RLock()
			f.groupQueue = append(f.groupQueue, arg[0].Interface())
			if len(f.groupQueue) >= f.group {
				arg = buildHandlerArg(f.groupQueue)
				f.groupQueue = make([]interface{}, 0)
				f.mtx.RUnlock()
				select {
				case f.outQueue <- arg:
				default:
					fmt.Println("Overflowing: The element disappears.", arg)
				}
			} else {
				f.mtx.RUnlock()
			}
		}
	}()

	return nil
}

func (f *groupFlow) Via(flow Flow) {
	flow.attachInBound(&In{
		queue: f.outQueue,
		run:   f.run,
	})
}

func (f *groupFlow) To(sink Sink) {
	sink.attachInBound(&In{
		queue: f.outQueue,
		run:   f.run,
	})
}
