package application

import (
	"errors"
	"github.com/sabouaram/reporting-reader/internal/config"
	gmail2 "github.com/sabouaram/reporting-reader/internal/domain/repositories/gmail"
	"log"
)

type Application struct {
	Config *config.Config
	gmail  *gmail2.GmailRepository
}

func NewApplication(config *config.Config, gmailrepo *gmail2.GmailRepository) (a *Application, err error) {
	if config.Username != "" && config.AuthorizedHTTPClient != nil {
		if err != nil {
			return nil, errors.New("Error Failed to create a new Application instance")
		}
		log.Println("Authorized Gmail HTTP Client Service Created Successfully => Connection Established")
		return &Application{
			Config: config,
			gmail:  gmailrepo,
		}, nil
	} else {
		return nil, errors.New("Failed to create new Application instance")
	}
}

// Starts the application
func (a *Application) StartApp() (err error) {
	// Listing User messages based on the defined Filter query
	listRes, err := a.gmail.ListMessages()
	if err != nil {
		// No Messages correspond to the Query :-)
		return err
	}
	// Filling Gmail repo instance with the received filtered list messages :-)
	err = a.gmail.SetMessages(listRes)
	if err != nil {
		// No Messages correspond to the filter query :-)
		return err
	}
	// Mark Processed mails as readed :-)
	err = a.gmail.MarkAsReaded()
	if err != nil {
		// Failed to mark a mail/mails as readed :-)
		return err
	}
	attmap, err := a.gmail.GetAttachmentsIds()
	if err != nil {
		// Error in getting attchments IDs :-)
		return err
	}
	log.Println("=> Processing ", len(attmap), "Attachments")
	err = a.gmail.GetAttachments(attmap)
	if err != nil {
		// Error in getting attchments IDs :-)
		return err
	}
	// Reset application messages slice
	a.gmail.ResetMessages()
	return nil
}
