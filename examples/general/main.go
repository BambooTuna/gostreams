package main

import (
	"fmt"
	"gostreams/stream"
	"sync"
	"time"
)

func main() {
	source := stream.NewSource()
	flow1 := stream.NewMapFlow(func(item string) string {
		return item + " --->"
	})
	flow2 := stream.NewMapFlow(func(item string) string {
		return "<--- " + item
	})
	groupFlow := stream.NewGroupFlow(5)
	sink := stream.NewPlainSink(func(item []interface{}) {
		//time.Sleep(time.Millisecond * 100)
		fmt.Println("Sink", item)
	})

	source.Via(flow1)
	flow1.Via(flow2)
	flow2.Via(groupFlow)
	groupFlow.To(sink)

	go func() {
		i := 0
		for {
			time.Sleep(time.Millisecond * 100)
			v := fmt.Sprintf("item %d", i)
			source.Publish(v)
			i++
		}
	}()

	err := sink.Run()
	fmt.Println(err)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	fmt.Println("end")
}
