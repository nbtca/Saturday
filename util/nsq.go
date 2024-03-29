package util

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nsqio/go-nsq"
)

var NSQProducer *nsq.Producer
var LogTopic = os.Getenv("LOG_TOPIC")
var EventTopic = os.Getenv("EVENT_TOPIC")

func GetNSQProducer() *nsq.Producer {
	nsqHost := os.Getenv("NSQ_HOST")
	nsqConfig := nsq.NewConfig()
	nsqConfig.AuthSecret = os.Getenv("NSQ_SECRET")

	if nsqHost != "" {
		var err error
		NSQProducer, err = nsq.NewProducer(nsqHost, nsqConfig)
		if err != nil {
			getLogger().Fatal(err)
		}
	}

	return NSQProducer
}
