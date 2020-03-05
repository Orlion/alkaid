package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Srv  *http.Server
	conf *Conf
}

type Conf struct {
	Debug        int8
	ReadTimeout  int64
	WriteTimeout int64
	Addr         string
}

func New(conf *Conf) *Server {
	if conf.Debug != 0 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &http.Server{
		Handler:      gin.New(),
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
	}

	return &Server{
		Srv:  srv,
		conf: conf,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		err = fmt.Errorf("server init error: [%w]", err)
		return err
	}

	s.Srv.Serve(listener)

	return nil
}

func (s *Server) Exit() {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	s.Srv.Shutdown(ctx)
}
