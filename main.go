package main

import (
	"log"
	"os"

	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/config"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/service"
)

func init() {
	if err := config.LoadEnv(); err != nil {
		if environment := os.Getenv("ENV"); environment == "developement" {
			log.Fatalf("fatal error config file: %s", err)
		}
	}
}

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("fatal error connect database: %s", err)
	}

	config.ConnectMailer(
		os.Getenv("MAILER_HOST"),
		os.Getenv("MAILER_USERNAME"),
		os.Getenv("MAILER_PASSWORD"),
	)

	q := queries.NewQueries(db)
	s := service.NewService()

	server := api.NewServer(q, s)

	server.Start(os.Getenv("PORT"))
}
