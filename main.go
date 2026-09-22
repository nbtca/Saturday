package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/router"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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
		if err := readConsulConfig(consulAddr, consulKey); err != nil {
			return fmt.Errorf("failed at reading config from consul: %w", err)
		}
		go func() {
			for {
				time.Sleep(time.Second * 5)
				if err := readConsulConfig(consulAddr, consulKey); err != nil {
					util.Logger.Errorf("unable to read remote config: %v", err)
				}
			}
		}()
	}
	return nil
}

var consulClient = &http.Client{Timeout: 10 * time.Second}

func readConsulConfig(addr, key string) error {
	if !strings.Contains(addr, "://") {
		addr = "http://" + addr
	}
	req, err := http.NewRequest(http.MethodGet, addr+"/v1/kv/"+key+"?raw", nil)
	if err != nil {
		return err
	}
	if token := viper.GetString("CONSUL_HTTP_TOKEN"); token != "" {
		req.Header.Set("X-Consul-Token", token)
	}
	res, err := consulClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("consul returned %s", res.Status)
	}
	viper.SetConfigType("json")
	return viper.ReadConfig(res.Body)
}

func main() {

	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	util.InitValidator()
	util.InitDialer()
	util.InitGithubClient()

	repo.InitDB()
	defer repo.CloseDB()

	service.LogtoServiceApp = service.MakeLogtoService(viper.GetString("logto.endpoint"))
	util.Logger.Debug("LogtoService initialized with endpoint: " + viper.GetString("logto.endpoint"))

	r := router.SetupRouter()

	viper.SetDefault("port", 4000)
	port := viper.GetInt("port")

	util.Logger.Infof("Starting server at %d...", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
