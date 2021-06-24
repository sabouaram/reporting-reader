package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"github.com/sabouaram/reporting-reader/config"
	"github.com/360EntSecGroup-Skylar/excelize"
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
				switch ext {
				case "csv":
					records := csvReader(decoded)
					log.Println(records)
				case "xlsx":
					xlsxReader(decoded)
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

// Processing Bliink csv reports
func csvReader(data []byte) (records []string) {
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
	return records
}

// Processing Bliink xlsx reports
func xlsxReader(data []byte) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		log.Println("Failed to convert received bytes to excelize file pointer ")
	}
	sheetMap := f.GetSheetMap()
	for _ , v := range sheetMap {
		log.Println("SHEET: ", v)
		for _, row := range f.GetRows(v) {
			for _, colCell := range row {
				log.Println(colCell)
			}
		}

	}


}
