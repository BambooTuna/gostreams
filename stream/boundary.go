package stream

type Item string

type In struct {
	queue chan Item
	run   func() error
}

type Out struct {
	queue chan Item
	run   func() error
}
