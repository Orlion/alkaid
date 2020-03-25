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
	return ek.Ekeeper.AddQueue(queueName, queue, limit)
}

func (ek *Ek) SyncSubscribe(eventName string, listener event.Listener) {
	return ek.Ekeeper.SyncSubscribe(eventName, listener)
}

func (ek *Ek) AsyncSubscribe(eventName string, listener event.Listener) {
	return ek.Ekeeper.AsyncSubscribe(eventName, listener)
}

func (ek *Ek) Publish(e *Event, queueName string) (c handleCode, err error) {
	return ek.Ekeeper.Publish(e, queueName)
}

func (ek *Ek) Listen() {
	return ek.Ekeeper.Listen()
}

func (ek *Ek) Exit() {
	return ek.Ekeeper.Exit()
}
