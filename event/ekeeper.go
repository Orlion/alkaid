package event

import (
	"errors"
	"sync"
	"time"

	"github.com/Orlion/alkaid/client"
	"github.com/sirupsen/logrus"
)

type Ekeeper struct {
	d      *dispatcher
	queues map[string]*queuePair
	isExit bool
	logger *client.Log
}

func NewEkeeper(logger *client.Log) (er *Ekeeper, err error) {
	er = &Ekeeper{
		d:      newDispatcher(),
		queues: make(map[string]*queuePair),
		logger: logger,
	}

	return
}

// 添加可用队列
func (ek *Ekeeper) AddQueue(queueName string, queue Queue, limit int32) {
	if limit < 1 {
		limit = 1
	}

	ek.queues[queueName] = &queuePair{
		queue:   queue,
		limiter: make(chan struct{}, limit),
	}
}

// 订阅同步事件
func (ek *Ekeeper) SyncSubscribe(eventName string, listener Listener) {
	ek.d.addSyncListener(eventName, listener)
}

// 订阅异步事件
func (ek *Ekeeper) AsyncSubscribe(eventName string, listener Listener) {
	ek.d.addAsyncListener(eventName, listener)
}

func (ek *Ekeeper) Publish(e *Event, queueName string) (c handleCode, err error) {
	var (
		exist bool
		qp    *queuePair
	)

	if ek.d.isExistAsyncListener(e) {
		// 写入到持久化队列中
		if qp, exist = ek.queues[queueName]; !exist {
			err = errors.New("Can not found queue:" + queueName)
			return
		}

		if err = qp.queue.Push(e); err != nil {
			return
		}
	}

	// 通知同步监听器处理事件
	ek.d.dispatch(e, false)

	return
}

func (ek *Ekeeper) Listen() {
	ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper Run...")
	waitGroup := sync.WaitGroup{}
	for _, qPair := range ek.queues {
		ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper Listen queue listen")
		waitGroup.Add(1)
		go func(qPair *queuePair) {
		I:
			for {
				time.Sleep(1 * time.Second)

				events, err := qPair.queue.Pull()
				if err != nil {
					ek.logger.Error(logrus.Fields{
						"err": err,
					}, "[App] Ekeeper Listen queue pull err")
					continue // TODO
				}

				for _, e := range events {
					ek.logger.Trace(logrus.Fields{
						"e": e,
					}, "[App] Ekeeper Listen event")

					if ek.isExit {
						ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper Listen receive exit")
						break I
					}

					<-qPair.limiter
					waitGroup.Add(1)
					go func(e *Event) {
						ek.d.dispatch(e, true)
						qPair.queue.Ack(e.Id)
						qPair.limiter <- struct{}{}
						waitGroup.Done()
					}(e)
				}
			}
			waitGroup.Done()
		}(qPair)
	}
	waitGroup.Wait()
	ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper exit...")
}

func (ek *Ekeeper) Exit() {
	ek.isExit = true
}
