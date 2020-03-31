package event

import "time"

// 事件抽象
type Event struct {
	Id         int64
	Name       string
	Data       map[string]interface{}
	CreateTime time.Time
}
