package app

import (
	"context"
	"encoding/base64"
	"github.com/sabouaram/reporting-reader/config"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"strings"
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
			mRes := getMessageAttachments(msg)
			for attID, ext := range mRes {
				Attachement, err := gmailService.Users.Messages.Attachments.Get(a.Config.Username, v.Id, attID).Do()
				if err != nil {
					log.Println("Error getting the attachment %v of Message %v", attID, v.Id)
				}
				decoded, err := base64.URLEncoding.DecodeString(Attachement.Data)
				if err != nil {
					log.Println("ERROR IN DECODING ATTACHMENT")
				}
				// Just for testing
				err = ioutil.WriteFile("Attachment."+ext, decoded, 0644)
				if err != nil {
					log.Printf("Failed to write attachment")
				}

			}

		}

	}

}

// Get AttachmentIDs of a single Message Body that contains CSV or XLSX files
// return a map: AttId=>extension
func getMessageAttachments(message *gmail.Message) (MapAttIdExt map[string]string) {
	mRes := make(map[string]string)
	var parts = message.Payload.Parts
	for _, v := range parts {
		if v.Filename != "" && len(v.Filename) > 0 && checkType(v.Filename) == true {
			mRes[v.Body.AttachmentId] = getType(v.Filename)
			log.Println(v.Filename)
		}
	}
	return mRes
}

// Return File extension
func getType(Filename string) string {
	s := strings.Split(Filename, ".")
	return s[len(s)-1]
}

// Check if the attachment file is a csv or an xlsx
func checkType(Filename string) bool {
	if strings.Contains(Filename, ".xlsx") == true || strings.Contains(Filename, ".csv") == true {
		return true
	} else {
		return false
	}
}
