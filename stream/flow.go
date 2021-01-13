package stream

type Flow interface {
	Via(flow Flow)
	To(sink Sink)
	run() error
	attachInBound(in *In)
}
