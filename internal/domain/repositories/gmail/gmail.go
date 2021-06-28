package gmailrepo

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/sabouaram/reporting-reader/internal/domain/repositories/utils"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

type GmailRepository struct {
	service  *gmail.Service
	messages []*gmail.Message
	Username string
	Query    string
}

// Returns a new gmail repo instance
func NewGmailRepository(username , query string, httpClient *http.Client) (gRepo *GmailRepository, err error) {
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, errors.New("Error Failed to create a new Gmail API Service")
	}
	log.Println("Gmail Backend Authorized HTTP Client Service Created Successfully => Connection Established")
	return &GmailRepository{
		service:  service,
		Username: username,
		Query:    query,
	}, nil
}

// Returns an RPC response => ListMessagesResponse that contains all mails based on GmailRepo's embedded Filter query
func (gRepo *GmailRepository) ListMessages() (*gmail.ListMessagesResponse, error) {
	q := gRepo.service.Users.Messages.List(gRepo.Username).Q(gRepo.Query)
	list, err := q.Do()
	if err != nil {
		return nil, errors.New("No messages correspond to the defined query ")
	} else {
		return list, nil
	}
}

// Fills GmailRepo's gmail Messages slice
func (gRepo *GmailRepository) SetMessages(list *gmail.ListMessagesResponse) error {
	if len(list.Messages) > 0 {
		for _, v := range list.Messages {
			msg, err := gRepo.service.Users.Messages.Get(gRepo.Username, v.Id).Do()
			if err != nil {
				return errors.New("Failed to get an existed user message")
			}
			gRepo.messages = append(gRepo.messages, msg)
		}
		return nil
	} else {
		return errors.New("There is no mails that corresponds to the application's query => Unable to fill application's messages from an empty list response")
	}
}

// Label as Read all the processed mails
func (gRepo *GmailRepository) MarkAsReaded() error {
	if len(gRepo.messages) > 0 {
		for _, v := range gRepo.messages {
			_, err := gRepo.service.Users.Messages.Modify(gRepo.Username, v.Id, &gmail.ModifyMessageRequest{
				RemoveLabelIds: []string{"UNREAD"},
			}).Do()
			if err != nil {
				return errors.New("Failed to mark a processed email as readed")
			}
		}
		return nil
	} else {
		return errors.New("Unable to mark as readed an empty mail list ")
	}
}

// Returns a map of mails Attachments IDs and files extensions
// Attachment Id => file extension
func (gRepo *GmailRepository) GetAttachmentsIds() (mapAttIdExt map[string]string, err error) {
	if len(gRepo.messages) > 0 {
		mRes := make(map[string]string)
		for _, msg := range gRepo.messages {
			if len(msg.Payload.Parts) > 0 {
				for _, v := range msg.Payload.Parts {
					if v.Filename != "" && len(v.Filename) > 0 && utils.CheckType(v.Filename) == true {
						mRes[v.Body.AttachmentId] = utils.GetType(v.Filename)
						log.Println(v.Filename)
					}
				}
			} else {
				return nil, errors.New("Unable to get Attachments IDs of an emptu payload")
			}
		}
		return mRes, nil
	} else {
		return nil, errors.New("Unable to get attachments IDs & extensions of an empty application's messages slice")
	}
}

// Reset Application messages slice
func (gRepo *GmailRepository) ResetMessages() {
	gRepo.messages = []*gmail.Message{}
}

// Get attachments data, for now we're just printing the content
func (gRepo *GmailRepository) GetAttachments(mapAttIdExt map[string]string) error {
	if len(gRepo.messages) > 0 {
		for _, v := range gRepo.messages {
			for attID, extension := range mapAttIdExt {
				Attachement, err := gRepo.service.Users.Messages.Attachments.Get(gRepo.Username, v.Id, attID).Do()
				if err != nil {
					return errors.New("Failed to get a message attachment")
				}
				decoded, err := base64.URLEncoding.DecodeString(Attachement.Data)
				if err != nil {
					return errors.New("Failed in decoding an attachment ")
				}
				switch extension {
				case "csv":
					records, err := utils.CsvReader(decoded)
					if err == nil {
						log.Println(records)
					}
				case "xlsx":
					colcells, err := utils.XlsxReader(decoded)
					if err == nil {
						log.Println(colcells)
					}
				}

			}
		}
		return nil
	} else {
		return errors.New("Unable to get attachments data of an empty application's messages slice")
	}
}
