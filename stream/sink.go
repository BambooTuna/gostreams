package stream

type Sink interface {
	attachInBound(in *In)
	Run() error
}
