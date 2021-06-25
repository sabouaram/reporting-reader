package main

import (
	"github.com/sabouaram/reporting-reader/filters"
	"log"

	"github.com/sabouaram/reporting-reader/app"
	"github.com/sabouaram/reporting-reader/config"
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
	app := &app.Application{
		Config: conf,
		Filter: f,
	}
	app.StartApp()
}
