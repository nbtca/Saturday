package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nbtca/saturday/container"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/util"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func initConfig() error {
	if err := godotenv.Load(); err != nil {
		util.Logger.Warnf("Error loading .env file: %v", err)
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
			return fmt.Errorf("failed at reading config from consul: %w", err)
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

	// Initialize dependency injection container
	container := container.NewContainer()
	util.Logger.Debug("Dependency injection container initialized")

	r := router.SetupRouter(container)

	viper.SetDefault("port", 4000)
	port := viper.GetInt("port")

	util.Logger.Infof("Starting server at %d...", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
