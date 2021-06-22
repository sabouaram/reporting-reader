package app

import (
	"log"
	"os"

	"github.com/sabouaram/reporting-reader/config"
	"github.com/sabouaram/reporting-reader/models"
	"google.golang.org/api/gmail/v1"
)

type Application struct {
	Config *config.Config
	Models *models.Models
}

func (*Application) StartApp(conf *config.Config, secret_file string) {
	// Create a new gmail service using the client
	if conf.AuthorizedHTTPClient != nil && secret_file != "" {
	gmailService, err := gmail.New(conf.AuthorizedHTTPClient)
		if err != nil {
			log.Printf("Error: %v", err)
		}

	} else {
		l
	}

}


func (*Application) ParseCSV(file_name string) {
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}	
	defer file.Close()
	reader := csv.
}



