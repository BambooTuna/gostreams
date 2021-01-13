package main

import (
	"fmt"
	"gostreams/stream"
	"sync"
	"time"
)

func main() {
	source := stream.NewSource()
	flow1 := stream.NewFlow(func(item string) string {
		return item + "!!!!"
	})
	flow2 := stream.NewFlow(func(item string) string {
		return item + "....."
	})
	sink := stream.NewSink(func(item string) {
		time.Sleep(time.Millisecond * 100)
		fmt.Println("Sink", item)
	})

	source.Via(flow1)
	flow1.Via(flow2)
	flow2.To(sink)

	go func() {
		i := 0
		for {
			time.Sleep(time.Millisecond * 1000)
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
