package main

import (
	"log"

	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/config"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/spf13/viper"
)

func init() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}
}

func main() {
	server := api.NewServer()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("fatal error connect database: %s", err)
	}

	server.Queries = queries.NewQueries(db)

	server.Start(viper.GetString("app.port"))
}
