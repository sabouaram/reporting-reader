package filters

import (
	"errors"
	"log"
	"strings"
	"time"
)

const layout = "2006/01/02"

type GmailFilter struct {
	Query string
}

/*
Generates Gmail search query to filter mails with a specified attachments extension string slice, sender, receiver and date NewFilter([]string{"xlsx","pdf","csv"}, "","salim@bliink.io","")
will generate a query for listing all mails that contains any of the specified extension addressed to salim@bliink.io in today's date
*/
func NewFilter(attextensions []string, sender string, receiver string, date string) (*GmailFilter, error) {
	if len(attextensions) > 0 && receiver != "" {
		if date == "" {
			date = time.Now().Format(layout)
		} else {
			dt, err := time.Parse(layout, date)
			if err != nil {
				log.Println("Invalid Date Format entry")
				return nil, errors.New("Invalid Date Format entry")
			}
			date = dt.Format(layout)
		}
		return &GmailFilter{
			Query: To + receiver + " " + From + sender + " " + AfterDate + date + HasAttach + strings.Join(attextensions[:], " OR "),
		}, nil
	} else {
		return nil, errors.New("Unspecified Attachments extension and receiver user")
	}
}
