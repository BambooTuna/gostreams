package stream

import (
	"fmt"
	"reflect"
)

type In struct {
	queue chan []reflect.Value
	run   func() error
}

func isValidHandler(fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	return nil
}

func buildHandlerArg(arg interface{}) []reflect.Value {
	return []reflect.Value{reflect.ValueOf(arg)}
}
