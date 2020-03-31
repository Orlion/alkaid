package event

type Queue interface {
	Push(*Event) error
	Pull() ([]*Event, error)
	Ack(id int64) error
}

type queuePair struct {
	queue   Queue
	limiter chan struct{}
}
