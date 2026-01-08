package config

import (
	"fmt"
	"github.com/Out-Of-India-Theory/oit-go-commons/config"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var configuration *Configuration

type Configuration struct {
	AppConfig            AppConfig
	ServerConfig         config.AppConfig
	DatabaseConfig       config.PostgresConfig
	OMSClientConfig      HttpClientConfig
	PlatformClientConfig HttpClientConfig
}

type HttpClientConfig struct {
	Address string
	Timeout time.Duration
	ApiKey  string
}

type AppConfig struct {
	DefaultSignedURLTTL time.Duration
}

func addConfigPath(v *viper.Viper) {
	v.AddConfigPath(".")
	v.AddConfigPath("config")
}

func init() {

	if os.Getenv("TEST_MODE") == "true" {
		fmt.Println("Skipping config loading in test mode")
		configuration = &Configuration{}
		return
	}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	addConfigPath(v)
	v.AutomaticEnv()
	var err error
	if err = v.ReadInConfig(); err != nil {
		fmt.Printf("error while reading config file, %v\n", err)
		panic(err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err = v.Unmarshal(&configuration); err != nil {
		fmt.Printf("error while deserializing config, %v\n", err)
		panic(err)
	}
}

func GetConfig() *Configuration {
	return configuration
}
