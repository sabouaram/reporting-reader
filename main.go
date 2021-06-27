package main

import (
	"github.com/sabouaram/reporting-reader/config"
	"github.com/sabouaram/reporting-reader/filters"
	"github.com/sabouaram/reporting-reader/internal/application"
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
	app, err := application.NewApplication(conf, f)
	if err != nil {
		log.Println(err)
	}
	app.StartApp()
}
