package flog

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Logger interface {
		Infof(format string, v ...interface{})
		Errorf(format string, v ...interface{})
		Debugf(format string, v ...interface{})
		Warnf(format string, v ...interface{})
		Panicf(format string, v ...interface{})
		SetLocal(pkgName string)
		Close() error
	}
	Flog struct {
		Logger *logrus.Entry
		Writer io.WriteCloser
	}
)

func New() *Flog {
	l := &Flog{}

	logger := logrus.New()
	logger.SetFormatter(getFormatter())
	logger.SetLevel(getLevel())

	out := getOutPut()
	logger.SetOutput(out)

	l.Writer = out
	l.Logger = logrus.NewEntry(logger)

	return l
}
func getFormatter() logrus.Formatter {
	//by default, our logger use text format
	// we will implement get format base one env file in latter
	var formatter logrus.Formatter
	formatter = &logrus.TextFormatter{
		TimestampFormat: time.RFC822,
	}
	return formatter
}
func getLevel() logrus.Level {
	// we will implement get level base one env file in latter
	return logrus.DebugLevel
}
func getOutPut() io.WriteCloser {
	// we will implement get output base one env file in latter
	return os.Stdout
}
func (l Flog) Infof(format string, v ...interface{}) {
	l.Logger.Infof(format, v...)
}
func (l *Flog) SetLocal(pkgName string) {
	l.Logger = l.Logger.WithField("package", pkgName)
}
func (l Flog) Errorf(format string, v ...interface{}) {
	l.Logger.Errorf(format, v...)
}
func (l Flog) Debugf(format string, v ...interface{}) {
	l.Logger.Debugf(format, v...)
}
func (l Flog) Warnf(format string, v ...interface{}) {
	l.Logger.Warnf(format, v...)
}
func (l Flog) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}
func (l *Flog) Close() error {
	if l.Writer != nil {
		return l.Writer.Close()
	}
	return nil
}
