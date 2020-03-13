package http

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Orlion/alkaid/client"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Gin  *gin.Engine
	srv  *http.Server
	conf *Conf
	log  *client.Log
}

type Conf struct {
	Debug        int8
	ReadTimeout  int64
	WriteTimeout int64
	Addr         string
}

func New(conf *Conf, log *client.Log) (*Server, error) {
	if conf.Debug != 0 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	// disabled default writer
	gin.DefaultWriter = bufio.NewWriter(src)
	gin := gin.New()

	gin.Use(logrusLogger(log), recovery(log))

	srv := &http.Server{
		Handler:      gin,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
	}

	return &Server{
		srv:  srv,
		conf: conf,
		Gin:  gin,
		log:  log,
	}, nil
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		err = fmt.Errorf("server init error: [%w]", err)
		return err
	}

	s.log.Trace(logrus.Fields{
		"addr": s.conf.Addr,
	}, "[App] Http Run...")
	s.srv.Serve(listener)

	return nil
}

func (s *Server) Exit() {
	s.log.Trace(logrus.Fields{}, "[App] Http Exit...")

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	s.srv.Shutdown(ctx)
}
