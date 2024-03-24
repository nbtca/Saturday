package util

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

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
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	mw := io.MultiWriter(os.Stdout, src)
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
