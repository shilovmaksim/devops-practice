package logger

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

const (
	defaultLogPath     = "log/server.log"
	defaultServiceName = "unknown"
)

type Logger struct {
	*logrus.Entry
}

func NewTestLogger() *Logger {
	testLogger, _ := test.NewNullLogger()
	logEntry := testLogger.WithFields(logrus.Fields{"service": "unknown"})
	return &Logger{
		logEntry,
	}
}

func New(logPath, serviceName string, logLevel logrus.Level) *Logger {
	if logPath == "" {
		logPath = defaultLogPath
	}
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logLevel)

	file := newFileWriter(logPath)
	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)
	logEntry := logger.WithFields(logrus.Fields{"service": serviceName})

	return &Logger{
		logEntry,
	}
}

func newFileWriter(filePath string) *os.File {
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	return f
}
