package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"github.com/sabouaram/reporting-reader/filters"
	"io"
	"log"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/sabouaram/reporting-reader/config"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Application's struct embeds Config, GmailFilter & private gmail gRPC client instances
type Application struct {
	Config   *config.Config
	Filter   *filters.GmailFilter
	service  *gmail.Service
	messages []*gmail.Message
}

// Creates Application's new client gRPC service instance
func (a *Application) newGmailService() (err error) {
	a.service, err = gmail.NewService(context.Background(), option.WithHTTPClient(a.Config.AuthorizedHTTPClient))
	if err != nil {
		return errors.New("Error Failed to create a new Gmail API Service")
	}
	log.Println("Gmail Backend Authorized HTTP Client Service Created Successfully")
	return nil
}

// Application's constructor
func NewApplication(conf *config.Config, filter *filters.GmailFilter) (*Application, error) {
	if conf.Username != "" && conf.AuthorizedHTTPClient != nil && filter.Query != "" {
		return &Application{
			Config: conf,
			Filter: filter,
		}, nil
	} else {
		return nil, errors.New("Failed to create new Application instance")
	}

}

// Returns an RPC response => ListMessagesResponse that contains all mails based on Application's embedded Filter query
func (a *Application) listMessages() (*gmail.ListMessagesResponse, error) {
	query := a.service.Users.Messages.List(a.Config.Username).Q(a.Filter.Query)
	list, err := query.Do()
	if err != nil {
		return nil, errors.New("No messages correspond to the defined query ")
	} else {
		return list, nil
	}
}

// Fills Application's gmail Messages slice
func (a *Application) setMessages(list *gmail.ListMessagesResponse) error {
	if len(list.Messages) > 0 {
		for _, v := range list.Messages {
			msg, err := a.service.Users.Messages.Get(a.Config.Username, v.Id).Do()
			if err != nil {
				return errors.New("Failed to get an existed user message")
			}
			a.messages = append(a.messages, msg)
		}
		return nil
	} else {
		return errors.New("Unable to fill application's messages from an empty list")
	}
}

// Label as Read all the processed mails
func (a *Application) markAsReaded() error {
	if len(a.messages) > 0 {
		for _, v := range a.messages {
			_, err := a.service.Users.Messages.Modify(a.Config.Username, v.Id, &gmail.ModifyMessageRequest{
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
func (a *Application) getAttachmentsIds() (mapAttIdExt map[string]string, err error) {
	if len(a.messages) > 0 {
		mRes := make(map[string]string)
		for _, msg := range a.messages {
			if len(msg.Payload.Parts) > 0 {
				for _, v := range msg.Payload.Parts {
					if v.Filename != "" && len(v.Filename) > 0 && checkType(v.Filename) == true {
						mRes[v.Body.AttachmentId] = getType(v.Filename)
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

func (a *Application) getAttachments(mapAttIdExt map[string]string) error {
	if len(a.messages) > 0 {
		for _, v := range a.messages {
			for attID, extension := range mapAttIdExt {
				Attachement, err := a.service.Users.Messages.Attachments.Get(a.Config.Username, v.Id, attID).Do()
				if err != nil {
					return errors.New("Failed to get a message attachment")
				}
				decoded, err := base64.URLEncoding.DecodeString(Attachement.Data)
				if err != nil {
					return errors.New("Failed in decoding an attachment ")
				}
				switch extension {
				case "csv":
					records, err := csvReader(decoded)
					if err == nil {
						log.Println(records)
					}
				case "xlsx":
					colcells, err := xlsxReader(decoded)
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

//
func (a *Application) StartApp() {
	// Create a new gmail service
	err := a.newGmailService()
	if err != nil {
		// Failed to create new gmail API client grpc service :-)
		log.Println(err)
		return
	}
	// Listing User messages based on the defined Filter query
	listRes, err := a.listMessages()
	if err != nil {
		// No Messages correspond to the Query :-)
		log.Println(err)
		return
	}
	// Filling Application instance with the received filtered list messages :-)
	err = a.setMessages(listRes)
	if err != nil {
		// No Messages correspond to the filter query :-)
		log.Println(err)
		return
	}
	log.Println("=> Processing ", len(a.messages), "mail")
	// Mark Processed mails as readed :-)
	err = a.markAsReaded()
	if err != nil {
		// Failed to mark a mail/mails as readed :-)
		log.Println(err)
		return
	}
	attmap, err := a.getAttachmentsIds()
	if err != nil {
		// Error in getting attchments IDs :-)
		log.Println(err)
		return
	}
	err = a.getAttachments(attmap)
	if err != nil {
		// Error in getting attchments IDs :-)
		log.Println(err)
		return
	}

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

// Processing Bliink csv reports
func csvReader(data []byte) (records []string, err error) {
	if len(data) > 0 {
		Data := string(data)
		tmp := ""
		r := csv.NewReader(strings.NewReader(Data))
		r.Comment = '#' // Comment symbol
		r.Comma = ','   // CSV Separator
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if len(record) > 0 {
				tmp = strings.Join(record, "")
				records = append(records, tmp)
			}
		}
		return records, nil
	}
	return nil, errors.New("Empty File bytes slice")
}

// Processing Bliink xlsx reports
func xlsxReader(data []byte) (colCells []string, err error) {
	if len(data) > 0 {
		f, err := excelize.OpenReader(bytes.NewReader(data))
		if err != nil {
			return nil, errors.New("Failed to convert received bytes to excelize file pointer ")
		}
		sheetMap := f.GetSheetMap()
		for k, v := range sheetMap {
			log.Println("SHEET", k, ":", v)
			rows, err := f.GetRows(v)
			if err != nil {
				return nil, errors.New("Failed in processing a row in xlsx file")
			}
			for _, row := range rows {
				for _, colcell := range row {
					colCells = append(colCells, colcell)
				}
			}
		}
		return colCells, nil
	}
	return nil, errors.New("Empty File bytes slice")
}
