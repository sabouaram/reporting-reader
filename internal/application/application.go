package application

import (
	"errors"
	"github.com/sabouaram/reporting-reader/config"
	"github.com/sabouaram/reporting-reader/domain/repositories"
	"github.com/sabouaram/reporting-reader/filters"
	"log"
)

type Application struct {
	Config *config.Config
	gmail  *repositories.GmailRepository
}

func NewApplication(config *config.Config, filter *filters.GmailFilter) (a *Application, err error) {
	if config.Username != "" && config.AuthorizedHTTPClient != nil && filter.Query != "" {
		gmailRepo, err := repositories.NewGmailRepository(config.Username, filter.Query, config.AuthorizedHTTPClient)
		if err != nil {
			return nil, errors.New("Error Failed to create a new Gmail Repo")
		}
		log.Println("Authorized Gmail HTTP Client Service Created Successfully => Connection Established")
		return &Application{
			Config: config,
			gmail:  gmailRepo,
		}, nil
	} else {
		return nil, errors.New("Failed to create new Application instance")
	}
}

// Starts the application
func (a *Application) StartApp() {
	// Listing User messages based on the defined Filter query
	listRes, err := a.gmail.ListMessages()
	if err != nil {
		// No Messages correspond to the Query :-)
		log.Println(err)
		return
	}
	// Filling Gmail repo instance with the received filtered list messages :-)
	err = a.gmail.SetMessages(listRes)
	if err != nil {
		// No Messages correspond to the filter query :-)
		log.Println(err)
		return
	}
	// Mark Processed mails as readed :-)
	err = a.gmail.MarkAsReaded()
	if err != nil {
		// Failed to mark a mail/mails as readed :-)
		log.Println(err)
		return
	}
	attmap, err := a.gmail.GetAttachmentsIds()
	if err != nil {
		// Error in getting attchments IDs :-)
		log.Println(err)
		return
	}
	log.Println("=> Processing ", len(attmap), "Attachments")
	err = a.gmail.GetAttachments(attmap)
	if err != nil {
		// Error in getting attchments IDs :-)
		log.Println(err)
		return
	}
	// Reset application messages slice
	a.gmail.ResetMessages()

}
