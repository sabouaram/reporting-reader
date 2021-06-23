package app

import (
	"context"
	"encoding/base64"
	"io/ioutil"
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
		req := gmailService.Users.Messages.List("me").Q("from: :is:unread Has:attachment")
		r, err := req.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}
		log.Printf("Processing %v Unreaded mail with attachment...\n", len(r.Messages))
		// Processing Unreaded mails with attachement
		for _, v := range r.Messages {
			msg, err := gmailService.Users.Messages.Get(a.Config.Username, v.Id).Do()
			if err != nil {
				log.Fatal("Error while getting Unreaded messages")
			}
			msgAttachIds := getMessageAttachments(msg)

			for _, attID := range msgAttachIds {
				Attachement, err := gmailService.Users.Messages.Attachments.Get(a.Config.Username, v.Id, attID).Do()
				if err != nil {
					log.Println("Error getting the attachment %v of Message %v", attID, v.Id)
				}
				decoded, err := base64.URLEncoding.DecodeString(Attachement.Data)
				if err != nil {
					log.Println("ERROR IN DECODING ATTACHMENT")
				}
				// Just for testing
				err = ioutil.WriteFile("Attachment", decoded, 0644)
				if err != nil {
					log.Printf("Failed to write attachment")
				}

			}

		}

	}

}

/*
// Get AttachmentIDs of multiple messages
func getMessagesAttachments(messages []*gmail.Message) (attachIDS []string) {
	attachId := []string{}
	for _, v := range messages {
		attchmsg := getMessageAttachments(v)
		if len(attchmsg) > 0 {
			for k, _ := range attchmsg {
				attachId = append(attachId, attchmsg[k])
			}
		}
	}
	return attachId
}*/

// Get AttachmentIDs of a single Message Body
func getMessageAttachments(message *gmail.Message) (attachIDs []string) {
	var parts = message.Payload.Parts
	attachId := []string{}
	for _, v := range parts {
		if v.Filename != "" && len(v.Filename) > 0 {
			attachId = append(attachId, v.Body.AttachmentId)
		}
	}
	return attachId
}
