package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/sabouaram/reporting-reader/config"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io"
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
			log.Println("Error Failed to create a new Gmail API Service: %v", err)
			return
		}
		log.Println("Gmail Backend Authorized HTTP Client Service Created Successfully")
		// Listing User messages based on a defined Query
		req := gmailService.Users.Messages.List("me").Q("from: :is:unread Has:attachment")
		r, err := req.Do()
		if err != nil {
			// No Messages correspond to the Query
			log.Println("Unable to retrieve messages: %v", err)
			return
		}
		log.Println("=> Processing ", len(r.Messages), "unreaded mail with attachment")
		// Processing Unreaded mails with attachement
		for _, v := range r.Messages {
			msg, err := gmailService.Users.Messages.Get(a.Config.Username, v.Id).Do()
			if err != nil {
				log.Println("Error while getting Unreaded messages")
				return
			}
			mRes , err := getMessageAttachments(msg)
			// Mark the Processed mail as READ
			_ , err = gmailService.Users.Messages.Modify(a.Config.Username, v.Id, &gmail.ModifyMessageRequest{
				RemoveLabelIds: [] string{"UNREAD"},
			}).Do()
			log.Println(err)
			if err != nil {
				log.Println("Failed to get Messages Attachments IDs")
				return
			}
			for attID, ext := range mRes {
				Attachement, err := gmailService.Users.Messages.Attachments.Get(a.Config.Username, v.Id, attID).Do()
				if err != nil {
					log.Println("Error getting the attachment %v of Message %v", attID, v.Id)
				}
				decoded, err := base64.URLEncoding.DecodeString(Attachement.Data)
				if err != nil {
					log.Println("Error in decoding attachment %v of Message", attID, v.Id)
				}
				switch ext {
				case "csv":
					records , err := csvReader(decoded)
					if err == nil {
						log.Println(records)
					}
				case "xlsx":
					colcells , err := xlsxReader(decoded)
					if err == nil  {
						log.Println(colcells)
					}
				}

			}

		}

	}

}

// Get AttachmentIDs of a single Message Body that contains CSV or XLSX files
// return a map: AttId=>extension
func getMessageAttachments(message *gmail.Message) (MapAttIdExt map[string]string, err error ) {
	var parts = message.Payload.Parts
	if len(parts) > 0 {
		mRes := make(map[string]string)
		for _, v := range parts {
			if v.Filename != "" && len(v.Filename) > 0 && checkType(v.Filename) == true {
				mRes[v.Body.AttachmentId] = getType(v.Filename)
				log.Println(v.Filename)
			}
		}
		return mRes , nil
	}
	return nil, errors.New("Empty Message Parts")
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
		return records , nil
	}
	return nil , errors.New("Empty File bytes slice")
}

// Processing Bliink xlsx reports
func xlsxReader(data []byte) (colCells []string, err error){
	if len(data) > 0 {
		f, err := excelize.OpenReader(bytes.NewReader(data))
		if err != nil {
			return nil, errors.New("Failed to convert received bytes to excelize file pointer ")
		}
		sheetMap := f.GetSheetMap()
		for k , v := range sheetMap {
			log.Println("SHEET", k , ":", v)
			rows , err := f.GetRows(v)
			if err != nil {
				return nil, errors.New("Failed in processing a row in xlsx file")
			}
			for _ , row := range rows {
				for _ , colcell := range row {
					colCells = append(colCells , colcell)
				}
			}
		}
		return colCells , nil
	}
	return nil, errors.New("Empty File bytes slice")
}



