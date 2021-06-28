package main

import (
	"github.com/sabouaram/reporting-reader/internal/application"
	"github.com/sabouaram/reporting-reader/internal/config"
	gmailrepo "github.com/sabouaram/reporting-reader/internal/domain/repositories/gmail"
	"github.com/sabouaram/reporting-reader/internal/domain/repositories/gmail/filters"
	"log"
)

func main() {
	conf, err := config.NewConfig("Oauth_credentials.json", "salim@bliink.io")
	if err != nil {
		log.Println(err)
	}
	f, err := filters.NewFilter([]string{"xlsx", "csv"}, "", conf.Username, "")
	if err != nil {
		log.Println(err)
	}
	log.Println("Query:", f)
	repo , _ := gmailrepo.NewGmailRepository(conf.Username, f.Query, conf.AuthorizedHTTPClient)
	app, _ := application.NewApplication(conf, repo)
	if err != nil {
		log.Println(err)
	}
	app.StartApp()
}
