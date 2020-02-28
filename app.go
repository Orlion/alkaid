package alkaid

import (
	"flag"
	"fmt"

	"github.com/Orlion/alkaid/client"
	"github.com/Orlion/alkaid/config"
	"github.com/Orlion/alkaid/event"
)

type App struct {
	Clients *clients
	ekeeper *event.Ekeeper
}

type clients struct {
	Mysql *client.Mysql
	Log   *client.Log
}

func NewApp() (app *App, err error) {
	var (
		conf  *config.Conf
		mysql *client.Mysql
		log   *client.Log
	)

	flag.Parse()

	if conf, err = config.New(); err != nil {
		err = fmt.Errorf("配置错误: [%w]", err)
		return
	}

	if mysql, err = client.NewMysql(conf.Clients.Mysql); err != nil {
		err = fmt.Errorf("配置错误: [%w]", err)
		return
	}

	if log, err = client.NewLog(conf.Clients.Log); err != nil {
		err = fmt.Errorf("配置错误: [%w]", err)
		return
	}

	ekeeper := event.NewEkeeper()

	app = &App{
		Clients: &clients{
			Mysql: mysql,
			Log:   log,
		},
		ekeeper: ekeeper,
	}

	return
}
