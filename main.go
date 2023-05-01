package main

import (
	"log"
	"os"

	"github.com/kanatsanan6/hrm/api"
	"github.com/kanatsanan6/hrm/config"
	"github.com/kanatsanan6/hrm/queries"
)

// func init() {
// 	if err := config.LoadEnv(); err != nil {
// 		log.Fatalf("fatal error config file: %s", err)
// 	}
// }

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

	server := api.NewServer(q)

	server.Start(os.Getenv("PORT"))
}
