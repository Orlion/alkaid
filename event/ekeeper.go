package event

import (
	"errors"
	"sync"
	"time"
)

type Ekeeper struct {
	d      *dispatcher
	queues map[string]*queuePair
	isExit bool
}

func New() (er *Ekeeper) {
	er = &Ekeeper{
		d:      newDispatcher(),
		queues: make(map[string]*queuePair),
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
	waitGroup := sync.WaitGroup{}
	for _, qp := range ek.queues {
		waitGroup.Add(1)
		go func(qp *queuePair) {
			waitGroupChild := sync.WaitGroup{}
		I:
			for {
				time.Sleep(1 * time.Second)

				events, err := qp.queue.Pull()
				if err != nil {
					continue // TODO
				}

				for _, e := range events {
					if ek.isExit {
						break I
					}

					<-qp.limiter
					waitGroupChild.Add(1)
					go func(e *Event) {
						ek.d.dispatch(e, true)
						qp.queue.Ack(e.Id)
						qp.limiter <- struct{}{}
						waitGroupChild.Done()
					}(e)
				}
			}
			waitGroupChild.Wait()
			waitGroup.Done()
		}(qp)
	}
	waitGroup.Wait()
}

func (ek *Ekeeper) Exit() {
	ek.isExit = true
}
