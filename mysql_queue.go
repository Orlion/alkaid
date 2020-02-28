package alkaid

import (
	"time"

	"github.com/Orlion/alkaid/event"
	"github.com/jinzhu/gorm"
)

type MysqlQueue struct {
	conn      *gorm.DB
	tableName string
}

type EventModel struct {
	Id             int64     `gorm:"primary_key;column:id"`
	EventName      string    `gorm:"column:eventName"`
	EventData      string    `gorm:"column:eventData"`
	CreateDatetime time.Time `gorm:"column:createDatetime"`
}

func NewMysqlQueue(conn *gorm.DB, tableName string) *MysqlQueue {
	return &MysqlQueue{
		conn:      conn,
		tableName: tableName,
	}
}

func (q *MysqlQueue) Push(e event.Event) error {

	mysqlRow := &EventModel{}
	q.conn.Create(mysqlRow)
}

func (q *MysqlQueue) Pop() <-chan event.Event {

}
