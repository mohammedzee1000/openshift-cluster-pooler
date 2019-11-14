package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogHandler struct {
	logger *logrus.Entry
}

//NewLogger gets new logger logger
func NewLogger(name string, debug bool) *LogHandler {
	l := logrus.New()

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	l.SetReportCaller(debug)
	lh := LogHandler{logger: logrus.NewEntry(l).WithField("name", name)}
	return &lh
}

//Info logs an info log message
func (l LogHandler) Info(componentName string, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	e := l.logger.WithField("component-name", componentName)
	e.Infoln(msg)
}

//Error logs an error log message
func (l LogHandler) Error(componentName string, err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	e := l.logger.WithField("component-name", componentName).WithError(err)
	e.Errorln(msg)
}

//Fatal logs a fatal error message
func (l LogHandler) Fatal(componentName string, err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	e := l.logger.WithField("component-name", componentName).WithError(err)
	e.Fatalln(msg)
}
