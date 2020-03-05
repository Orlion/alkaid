package event

type dispatcher struct {
	syncListenersMap  map[string][]Listener // 同步监听器
	asyncListenersMap map[string][]Listener // 异步监听器
}

func newDispatcher() (d *dispatcher) {
	return &dispatcher{
		syncListenersMap:  make(map[string][]Listener),
		asyncListenersMap: make(map[string][]Listener),
	}
}

func (d *dispatcher) addSyncListener(eventName string, listener Listener) {
	d.syncListenersMap[eventName] = append(d.syncListenersMap[eventName], listener)
}

func (d *dispatcher) addAsyncListener(eventName string, listener Listener) {
	d.asyncListenersMap[eventName] = append(d.asyncListenersMap[eventName], listener)
}

// 是否存在异步监听器
func (d *dispatcher) isExistAsyncListener(e *Event) bool {
	_, exist := d.asyncListenersMap[e.Name]
	return exist
}

// 依次调用注册在事件上的监听器
func (d *dispatcher) dispatch(e *Event, async bool) (collector *HandleResCollector) {
	var (
		listeners []Listener
		exist     bool
	)

	collector = newHandleResCollector()

	if async {
		if listeners, exist = d.asyncListenersMap[e.Name]; !exist {
			return collector
		}
	} else {
		if listeners, exist = d.syncListenersMap[e.Name]; !exist {
			return collector
		}
	}

	for _, listener := range listeners {
		collector.collect(listener.Name(), listener.Handle(e))
	}

	return collector
}
