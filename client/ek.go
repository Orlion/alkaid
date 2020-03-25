package client

import "github.com/Orlion/alkaid/event"

type Ek struct {
	Ekeeper *event.Ekeeper
}

func NewEk(logger *Log) (ek *Ek) {
	ekeeper, err := event.NewEkeeper(logger)
	if err != nil {
		return
	}

	return &Ek{
		Ekeeper: ekeeper,
	}
}

func (ek *Ek) AddQueue(queueName string, queue event.Queue, limit int32) {
	ek.Ekeeper.AddQueue(queueName, queue, limit)
}

func (ek *Ek) SyncSubscribe(eventName string, listener event.Listener) {
	ek.Ekeeper.SyncSubscribe(eventName, listener)
}

func (ek *Ek) AsyncSubscribe(eventName string, listener event.Listener) {
	ek.Ekeeper.AsyncSubscribe(eventName, listener)
}

func (ek *Ek) Publish(e *event.Event, queueName string) (handleResList []*event.HandleRes, err error) {
	return ek.Ekeeper.Publish(e, queueName)
}

func (ek *Ek) Listen() {
	ek.Ekeeper.Listen()
}

func (ek *Ek) Exit() {
	ek.Ekeeper.Exit()
}
