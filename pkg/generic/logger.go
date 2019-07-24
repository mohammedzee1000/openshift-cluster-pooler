package generic

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func Init() {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,

	}
	logrus.SetFormatter(formatter)
}

type LogHander struct {
	logger *logrus.Entry
}

func NewLogger(name string) *LogHander {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	lh := LogHander{logger: logrus.NewEntry(l).WithField("name", name)}
	return &lh
}

func (l LogHander) Info(componentName string, format string, args ...interface{})  {
	msg := fmt.Sprintf(format, args...)
	e := l.logger.WithField("component-name", componentName)
	e.Infoln(msg)
}

func (l LogHander) Error(componentName string, err error,format string, args ...interface{})  {
	msg := fmt.Sprintf(format, args...)
	e := l.logger.WithField("component-name", componentName).WithError(err)
	e.Errorln(msg)
}

func (l LogHander) Fatal(componentName string, err error,format string, args ...interface{})  {
	msg := fmt.Sprintf(format, args...)
	e := l.logger.WithField("component-name", componentName).WithError(err)
	e.Fatalln(msg)
}