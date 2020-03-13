package client

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type LogConf struct {
	Dir      string
	FileName string
	Level    int32
}

type Log struct {
	logger *logrus.Logger
}

func NewLog(conf *LogConf) (log *Log, err error) {
	var (
		baseLogPath string
		writer      *rotatelogs.RotateLogs
	)

	log = &Log{}
	log.logger = logrus.New()

	log.logger.SetLevel(logrus.Level(conf.Level))

	baseLogPath = path.Join(conf.Dir, conf.FileName)
	if writer, err = rotatelogs.New(baseLogPath+".%Y%m%d%H", rotatelogs.WithLinkName(baseLogPath), rotatelogs.WithRotationTime(time.Hour), rotatelogs.WithMaxAge(time.Hour*24)); nil != err {
		return
	}

	pathMap := lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}

	log.logger.AddHook(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return
}

func (log *Log) Trace(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Trace(args)
}

func (log *Log) Debug(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Debug(args)
}

func (log *Log) Info(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Info(args)
}

func (log *Log) Warn(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Warn(args)
}

func (log *Log) Error(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Error(args)
}

func (log *Log) Fatal(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Fatal(args)
}

func (log *Log) Panic(fields logrus.Fields, args ...interface{}) {
	log.logger.WithFields(fields).Panic(args)
}
