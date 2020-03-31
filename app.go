package alkaid

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Orlion/alkaid/client"
	"github.com/sirupsen/logrus"
)

type App struct {
	name    string
	Clients *clients
}

type clients struct {
	Mysql *client.Mysql
	Log   *client.Log
	Ek    *client.Ek
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func New(name string) (app *App, err error) {
	var (
		conf  *Conf
		mysql *client.Mysql
		log   *client.Log
	)

	flag.Parse()

	if conf, err = newConf(); err != nil {
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

	app = &App{
		name: name,
		Clients: &clients{
			Mysql: mysql,
			Log:   log,
			Ek: client.NewEk(log),
		},
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
			app.Clients.Ek.Exit()
			app.Clients.Log.Trace(logrus.Fields{
				"signal": s,
			}, "[App] GraceExit begin")
			callback()
			// clients exit...
			app.Clients.Log.Trace(logrus.Fields{
				"signal": s,
			}, "[App] GraceExit end")
			return
		case syscall.SIGHUP:
		default:
		}
	}
}

func (app *App) GetName() string {
	return app.name
}
