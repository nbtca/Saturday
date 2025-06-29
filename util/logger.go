package util

import (
	"io"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type ContextLogger struct {
	*logrus.Logger
	Context *gin.Context
}

type NSQHookForError struct {
	Producer *nsq.Producer
}

// Fire 根据 logrus.Entry 发送消息到 NSQ
func (hook *NSQHookForError) Fire(entry *logrus.Entry) error {
	// 将日志消息发送到 NSQ
	byte, err := entry.Bytes()
	if err != nil {
		return err
	}
	return hook.Producer.Publish(LogTopic, byte)
}

// Levels 返回日志级别，这里返回 ErrorLevel,FatalLevel,PanicLevel
func (hook *NSQHookForError) Levels() []logrus.Level {
	levels := [3]logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
	return levels[:]
}

func getLogger() *logrus.Logger {
	//实例化
	logger := logrus.New()

	//写入到标准输出
	mw := io.MultiWriter(os.Stdout)
	logrus.SetOutput(mw)
	logger.Out = mw

	//设置日志级别
	logger.SetLevel(logrus.TraceLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if producer := GetNSQProducer(); producer != nil {
		logger.Hooks.Add(&NSQHookForError{
			Producer: producer,
		})
	}

	return logger
}

var Logger = getLogger()

// GetStackTrace returns the current stack trace as a string
func GetStackTrace(skip int) string {
	buf := make([]byte, 1024*8)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, len(buf)*2)
	}
}

// LogWithStackTrace logs an error with stack trace
func LogWithStackTrace(level logrus.Level, msg string, err error) {
	entry := Logger.WithField("stacktrace", GetStackTrace(2))
	if err != nil {
		entry = entry.WithError(err)
	}
	entry.Log(level, msg)
}

// LogErrorWithStackTrace logs an error with stack trace
func LogErrorWithStackTrace(msg string, err error) {
	LogWithStackTrace(logrus.ErrorLevel, msg, err)
}

// LogFatalWithStackTrace logs a fatal error with stack trace
func LogFatalWithStackTrace(msg string, err error) {
	LogWithStackTrace(logrus.FatalLevel, msg, err)
}
