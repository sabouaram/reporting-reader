package main

import (
	"log"

	"github.com/sabouaram/reporting-reader/app"
	"github.com/sabouaram/reporting-reader/config"
)

func main() {
	conf, err := config.NewConfig("Oauth_credentials.json", "salim@bliink.io")
	if err != nil {
		log.Println(err)
	}
	app := &app.Application{
		Config: conf,
	}
	app.StartApp()
}
