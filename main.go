package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/util"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func initConfig() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	consulAddr := viper.GetString("CONSUL_HTTP_ADDR")
	consulKey := viper.GetString("CONSUL_KEY")
	if consulAddr != "" {
		util.Logger.Debug("Using consul config", consulAddr)
		viper.AddRemoteProvider("consul", consulAddr, consulKey)
		viper.SetConfigType("json") // Need to explicitly set this to json
		err := viper.ReadRemoteConfig()
		if err != nil {
			return fmt.Errorf("failed at reading config: %w", err)
		}
		go func() {
			for {
				time.Sleep(time.Second * 5) // delay after each request

				// currently, only tested with etcd support
				err := viper.WatchRemoteConfig()
				if err != nil {
					util.Logger.Errorf("unable to read remote config: %v", err)
					continue
				}
			}
		}()
	}
	return nil
}

func main() {

	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	util.InitValidator()
	util.InitDialer()

	repo.InitDB()
	defer repo.CloseDB()

	r := router.SetupRouter()

	viper.SetDefault("port", 4000)
	port := viper.GetInt("port")

	r.Run(fmt.Sprintf(":%d", port))

	util.Logger.Infof("Starting server at %d...", port)
}
