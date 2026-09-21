package util

import (
	"log"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

var (
	nsqProducer   *nsq.Producer
	nsqProducerMu sync.Mutex
)

// Topics are read on use because package-level vars are initialized before main loads the config.
func LogTopic() string   { return viper.GetString("nsq.log_topic") }
func EventTopic() string { return viper.GetString("nsq.event_topic") }

func GetNSQProducer() *nsq.Producer {
	nsqProducerMu.Lock()
	defer nsqProducerMu.Unlock()
	if nsqProducer != nil {
		return nsqProducer
	}
	nsqHost := viper.GetString("nsq.host")
	if nsqHost == "" {
		return nil
	}
	nsqConfig := nsq.NewConfig()
	nsqConfig.AuthSecret = viper.GetString("nsq.secret")
	producer, err := nsq.NewProducer(nsqHost, nsqConfig)
	if err != nil {
		log.Printf("failed to create nsq producer: %v", err)
		return nil
	}
	nsqProducer = producer
	return nsqProducer
}
