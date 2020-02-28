package event

type Queue interface {
	Push(e Event) error
	Pop() <-chan Event
}

type queuePair struct {
	queue   Queue
	limiter chan struct{}
}
