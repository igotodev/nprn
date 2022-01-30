package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

var entry *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{entry}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

type writerHook struct {
	Writer   []io.Writer
	LogLevel []logrus.Level
}

func (wh *writerHook) Fire(entry *logrus.Entry) error {
	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range wh.Writer {
		_, err = w.Write([]byte(str))
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (wh *writerHook) Levels() []logrus.Level {
	return wh.LogLevel
}

// auto call if logger use in code
func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", fileName, frame.Line)

		},
		FullTimestamp: true,
		DisableColors: false,
	}

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Println(err)
	}

	logFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}

	l.SetOutput(io.Discard) // set default output is nothing

	// set custom output in logFile and os.Stdout
	l.AddHook(&writerHook{
		Writer:   []io.Writer{logFile, os.Stdout},
		LogLevel: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel) // max level

	entry = logrus.NewEntry(l)
}
