package stream

type Source interface {
	Publish(arg interface{})
	Via(flow Flow)
	To(sink Sink)
}
