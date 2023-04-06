package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server Server `mapstructure:",squash"`
}

type Server struct {
	Port string
}

func (c *Config) InitCfg() error {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&c)

	return err
}
