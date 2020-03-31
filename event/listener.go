package event

type Listener func(*Event) *HandleRes

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
