package alkaid

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Orlion/alkaid/client"
	"github.com/Orlion/alkaid/config"
	"github.com/Orlion/alkaid/event"
	"github.com/Orlion/alkaid/http"
	"github.com/sirupsen/logrus"
)

type App struct {
	context.Context
	Clients *clients
	Ekeeper *event.Ekeeper
	Server  *http.Server
}

type clients struct {
	Mysql *client.Mysql
	Log   *client.Log
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // 用满所有核心
}

func NewApp() (app *App, err error) {
	var (
		conf  *config.Conf
		mysql *client.Mysql
		log   *client.Log
	)

	flag.Parse()

	if conf, err = config.New(); err != nil {
		err = fmt.Errorf("config error: [%w]", err)
		return
	}

	if mysql, err = client.NewMysql(conf.Clients.Mysql); err != nil {
		err = fmt.Errorf("new mysql error: [%w]", err)
		return
	}

	if log, err = client.NewLog(conf.Clients.Log); err != nil {
		err = fmt.Errorf("new log error: [%w]", err)
		return
	}

	ekeeper := event.NewEkeeper(log)

	server, err := http.New(conf.Http, log)
	if err != nil {
		err = fmt.Errorf("new http error: [%w]", err)
		return
	}

	app = &App{
		Clients: &clients{
			Mysql: mysql,
			Log:   log,
		},
		Ekeeper: ekeeper,
		Server:  server,
	}

	return
}

func (app *App) GraceExit(callback func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			app.Clients.Log.Trace(logrus.Fields{
				"signal": s,
			}, "[App] GraceExit")
			app.Ekeeper.Exit()
			app.Server.Exit()
			callback()
			return
		case syscall.SIGHUP:
		default:
		}
	}
}
