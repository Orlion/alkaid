package event

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Ekeeper struct {
	d         *dispatcher
	queues    map[string]*queuePair
	isExit    bool
	waitGroup sync.WaitGroup
	logger    Log
}

func NewEkeeper(logger Log) (er *Ekeeper, err error) {
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

func (ek *Ekeeper) Publish(e *Event, queueName string) (handleResList []*HandleRes, err error) {
	var (
		exist bool
		qp    *queuePair
	)

	if ek.d.isExistAsyncListener(e) || queueName != "" {
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
	handleResList = ek.d.dispatch(e, false)

	return
}

func (ek *Ekeeper) Listen() {
	ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper Run...")
	ek.waitGroup = sync.WaitGroup{}
	for qName, qPair := range ek.queues {
		ek.logger.Trace(logrus.Fields{
			"queue": qName,
		}, "[App] Ekeeper queue listen...")
		ek.waitGroup.Add(1)
		go func(qp *queuePair) {
		I:
			for {
				if ek.isExit {
					ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper Listen receive exit")
					break I
				}

				time.Sleep(1 * time.Second)

				events, err := qp.queue.Pull()
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
						ek.logger.Trace(logrus.Fields{}, "[App] Ekeeper Listen receive exit1")
						break I
					}

					qp.limiter <- struct{}{}
					ek.waitGroup.Add(1)
					go func(e *Event) {
						ek.d.dispatch(e, true)
						qp.queue.Ack(e.Id)
						<-qp.limiter
						ek.waitGroup.Done()
					}(e)
				}
			}
			ek.waitGroup.Done()
		}(qPair)
	}
}

func (ek *Ekeeper) Exit() {
	ek.logger.Trace(logrus.Fields{}, "Ekeeper exit begin")
	ek.isExit = true
	ek.waitGroup.Wait()
}
