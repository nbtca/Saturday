package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/util"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {

	if err := godotenv.Load(); err != nil {
		util.Logger.Warning("Error loading .env file")
		log.Println(err)
	}

	viper.AutomaticEnv()
	viper.SetDefault("port", 4000)
	// https://github.com/bketelsen/crypt/blob/5cbc8cc4026c0c1d3bf9c5d4e5a30398f99c99a9/vendor/github.com/hashicorp/consul/api/api.go#L31
	consulAddr := viper.GetString("CONSUL_HTTP_ADDR")
	consulKey := viper.GetString("CONSUL_KEY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if consulAddr != "" {
		util.Logger.Debug("Using consul config", consulAddr)
		viper.AddRemoteProvider("consul", consulAddr, consulKey)
		viper.SetConfigType("json") // Need to explicitly set this to json
		err := viper.ReadRemoteConfig()
		if err != nil {
			util.Logger.Error("failed at reading config: ", err)
		}
		// log.Print(viper.GetString("hostname"))
		// log.Print(viper.GetString("port"))
	}

	util.InitValidator()
	util.InitDialer()

	repo.InitDB()
	defer repo.CloseDB()

	r := router.SetupRouter()

	port := viper.GetInt("port") // Will read from env var MYAPP_PORT if present
	if port == 0 {
		log.Fatal("$PORT must be set")
	}
	r.Run(fmt.Sprintf(":%d", port)) // listen and serve on

	util.Logger.Infof("Starting server at %d...", port)
}
