package event

import "github.com/sirupsen/logrus"

type Log interface {
	Trace(fields logrus.Fields, args ...interface{})
	Debug(fields logrus.Fields, args ...interface{})
	Info(fields logrus.Fields, args ...interface{})
	Warn(fields logrus.Fields, args ...interface{})
	Error(fields logrus.Fields, args ...interface{})
}
