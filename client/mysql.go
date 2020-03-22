package client

import (
	"github.com/pkg/errors"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct {
	conns map[string]*gorm.DB
}

type MysqlConf struct {
	Conns map[string]struct {
		Dsn             string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
		LogMode         bool
	}
}

func NewMysql(mysqlConf *MysqlConf) (mysql *Mysql, err error) {
	mysql = &Mysql{}
	mysql.conns = make(map[string]*gorm.DB)

	for key, connConf := range mysqlConf.Conns {
		if mysql.conns[key], err = gorm.Open("mysql", connConf.Dsn); nil != err {
			return
		}

		mysql.conns[key].DB().SetConnMaxLifetime(connConf.ConnMaxLifetime * time.Second)
		mysql.conns[key].DB().SetMaxIdleConns(connConf.MaxIdleConns)
		mysql.conns[key].DB().SetMaxOpenConns(connConf.MaxOpenConns)

		mysql.conns[key].LogMode(connConf.LogMode)
	}

	return
}

func (mysql *Mysql)Conn(name string) (conn *gorm.DB, err error) {
	conn, exist := mysql.conns[name]
	if !exist {
		err = errors.New("not found conn:" + name)
	}

	return
}