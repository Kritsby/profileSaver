package main

import (
	_ "dev/profileSaver/docs"
	"dev/profileSaver/internal/app"
	"dev/profileSaver/internal/config"
	"log"
)

var cfg config.Config

func init() {
	err := cfg.InitCfg()
	if err != nil {
		panic(err)
	}
}

// @title SHOP API
// @version 1.0
// @description API Server
// @BasePath /

// @securityDefinitions.basic BasicAuth
func main() {
	err := app.Run(cfg)
	if err != nil {
		log.Println(err)
	}
}
