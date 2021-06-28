package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
	"log"
	"strings"
)

// Return File extension
func GetType(Filename string) string {
	s := strings.Split(Filename, ".")
	return s[len(s)-1]
}

// Check if the attachment file is a csv or an xlsx
func CheckType(Filename string) bool {
	if strings.Contains(Filename, ".xlsx") == true || strings.Contains(Filename, ".csv") == true {
		return true
	} else {
		return false
	}
}

// Processing Bliink csv reports
func CsvReader(data []byte) (records []string, err error) {
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
				if strings.Contains(tmp, string(r.Comment)) == false {
					records = append(records, tmp)
				}
			}
		}
		return records, nil
	}
	return nil, errors.New("Empty File bytes slice")
}

// Processing Bliink xlsx reports
func XlsxReader(data []byte) (colCells []string, err error) {
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
