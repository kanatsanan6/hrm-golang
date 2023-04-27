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
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("fatal error connect database: %s", err)
	}

	q := queries.NewQueries(db)

	server := api.NewServer(q)

	server.Start(viper.GetString("app.port"))
}
