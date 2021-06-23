package app

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/sabouaram/reporting-reader/config"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Application struct {
	Config *config.Config
	//Models *models.Models
}

func (a *Application) StartApp() {
	// Create a new gmail service
	if a.Config.AuthorizedHTTPClient != nil {
		gmailService, err := gmail.NewService(context.Background(), option.WithHTTPClient(a.Config.AuthorizedHTTPClient))
		if err != nil {
			log.Println("Error: %v", err)
		}
		log.Printf("Gmail Backend Authorized HTTP Client Service Created Successfully")

		// Listing User messages based on a defined Query
		req := gmailService.Users.Messages.List("me").Q("from: :is:unread Has:attachment ")
		if err != nil {
			log.Println(err)
		}
		r, err := req.Do()

		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}
		log.Printf("Processing %v Unreaded mail with attachment...\n", len(r.Messages))

		// Processing Unreaded mails with attachement
		msgs := []*gmail.Message{}
		for _, v := range r.Messages {
			msg, err := gmailService.Users.Messages.Get("me", v.Id).Do()
			if err != nil {
				log.Fatal("Error while getting Unreaded messages")
			}
			log.Println(msg)
			msgs = append(msgs, msg)
		}
		for _, v := range msgs {
			attach, _ := gmailService.Users.Messages.Attachments.Get("me", v.Id, v.Payload.Body.AttachmentId).Do()
			decoded, err := base64.StdEncoding.DecodeString(attach.Data)
			if err == nil {
				log.Println("ATTACHMENT", decoded)
			}
		}

		/*decoded_attachments := []byte{}
		for _, v := range msgs {
			attach, _ := gmailService.Users.Messages.Attachments.Get("me", v.Id, v.Payload.Body.AttachmentId).Do()
			decoded, err := base64.URLEncoding.DecodeString(attach.Data)
			if err != nil {
				log.Println("Error", err)
			}
			decoded_attachments = append(decoded_attachments, decoded...)
		}*/

	}

}
