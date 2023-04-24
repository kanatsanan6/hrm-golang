package main

import (
	"log"

	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/config"
	"github.com/spf13/viper"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	server := api.NewServer()

	server.Start(viper.GetString("app.port"))
}
