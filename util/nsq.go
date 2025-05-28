package util

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

var NSQProducer *nsq.Producer
var LogTopic = viper.GetString("nsq.log_topic")
var EventTopic = viper.GetString("nsq.event_topic")

func GetNSQProducer() *nsq.Producer {
	nsqHost := viper.GetString("nsq.host")
	nsqConfig := nsq.NewConfig()
	nsqConfig.AuthSecret = viper.GetString("nsq.secret")

	if nsqHost != "" {
		var err error
		NSQProducer, err = nsq.NewProducer(nsqHost, nsqConfig)
		if err != nil {
			getLogger().Fatal(err)
		}
	}

	return NSQProducer
}
