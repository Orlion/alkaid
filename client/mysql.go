package client

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct {
	Conns map[string]*gorm.DB
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
	mysql.Conns = make(map[string]*gorm.DB)

	for key, connConf := range mysqlConf.Conns {
		if mysql.Conns[key], err = gorm.Open("mysql", connConf.Dsn); nil != err {
			return
		}

		mysql.Conns[key].DB().SetConnMaxLifetime(connConf.ConnMaxLifetime * time.Second)
		mysql.Conns[key].DB().SetMaxIdleConns(connConf.MaxIdleConns)
		mysql.Conns[key].DB().SetMaxOpenConns(connConf.MaxOpenConns)

		mysql.Conns[key].LogMode(connConf.LogMode)
	}

	return
}
