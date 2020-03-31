package alkaid

import (
	"errors"
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/Orlion/alkaid/client"
)

type Conf struct {
	Clients *client.Conf
}

var (
	confFilePath string
)

func init() {
	flag.StringVar(&confFilePath, "conf", "app.toml", "default config path")
}

func newConf() (conf *Conf, err error) {
	var (
		tmpConf *Conf
	)

	if confFilePath == "" {
		err = errors.New("config file path is required")
		return
	}

	if _, err = toml.DecodeFile(confFilePath, &tmpConf); nil != err {
		return
	}

	conf = tmpConf

	return
}
