package event

// 监听器抽象
type Listener interface {
	Handle(*Event) *HandleRes
	Name() string
}

type handleCode int8

const (
	HandleFailed  handleCode = iota - 1 // 失败
	HandleSuccess                       // 成功
	HandleError                         // 错误
)

type HandleRes struct {
	Code handleCode
	Msg  string
}

// 处理结果收集器
type HandleResCollector struct {
	listenerHandleResMap map[string]*HandleRes
}

func newHandleResCollector() *HandleResCollector {
	return &HandleResCollector{
		listenerHandleResMap: make(map[string]*HandleRes),
	}
}

func (collector *HandleResCollector) collect(listenerName string, res *HandleRes) {
	collector.listenerHandleResMap[listenerName] = res
}
