package models

import (
	"time"
)

type Models struct {
	Extension string
	Date      string
}

func NewModels(Extension string, Dateformat string) *Models {
	return &Models{
		Extension: Extension,
		Date:      time.Now().Format(Dateformat),
	}
}
