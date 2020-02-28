package event

import "errors"

type Ekeeper struct {
	d      *dispatcher
	queues map[string]queuePair
}

func NewEkeeper() (er *Ekeeper) {
	er = &Ekeeper{
		d:      newDispatcher(),
		queues: make(map[string]queuePair),
	}

	return
}

// 添加可用队列
func (ek *Ekeeper) AddQueue(queueName string, queue Queue, limit int32) {
	if limit < 1 {
		limit = 1
	}

	ek.queues[queueName] = queuePair{
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

func (ek *Ekeeper) Publish(e Event, queueName string) (c handleCode, err error) {
	var (
		exist     bool
		queuePair queuePair
	)

	if ek.d.isExistAsyncListener(e) {
		// 写入到持久化队列中
		if queuePair, exist = ek.queues[queueName]; !exist {
			err = errors.New("Can not found queue:" + queueName)
			return
		}

		if err = queuePair.queue.Push(e); err != nil {
			return
		}
	}

	// 通知同步监听器处理事件
	ek.d.dispatch(e, false)

	return
}

func (ek *Ekeeper) Listener() {
	for _, queuePair := range ek.queues {
		go func() {
			for e := range queuePair.queue.Pop() {
				<-queuePair.limiter
				go func() {
					ek.d.dispatch(e, true)
					queuePair.limiter <- struct{}{}
				}()
			}
		}()
	}
}
