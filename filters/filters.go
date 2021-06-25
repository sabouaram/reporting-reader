package filters

import (
	"time"
)

type Filter struct {
	Query string
	Date  string
}

func NewFilter(Query string, Dateformat string) *Filter {
	return &Models{
		Query: Query,
		Date:  time.Now().Format(Dateformat),
	}
}

// TO DO DYNAMIC QUERY CREATOR
