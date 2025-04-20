package util

import (
	"io"
	"os"

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
	logger.SetLevel(logrus.DebugLevel)

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
