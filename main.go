package main

import (
	"bytes"
	"strconv"
	"strings"
)

const ESCAPE_CHAR = '\\'
const QUOTE_CHAR = '"'
const DELIMITER = ','
const NEW_LINE = '\n'

func ParseRows(csv string) [][]any {
	openQuote := -1
	nextEscapedCharIdx := -1
	var rows [][]any = [][]any{}
	var row []any = []any{}
	var rowIdx = 0
	field := bytes.NewBufferString("")
	var fieldIdx = 0
	for i, c := range csv {
		switch c {
		case ESCAPE_CHAR:
			// if this escape char is not escaped
			if nextEscapedCharIdx != i {
				// flag the next char as escaped
				nextEscapedCharIdx = i + 1
			} else {
				// if this escape char is itself escaped, then store it
				field.WriteRune(c)
			}
		case QUOTE_CHAR:
			// if this char is not escaped
			if nextEscapedCharIdx != i {
				// if we are in quotes already...
				if openQuote > -1 {
					if openQuote+1 == i {
						// don't record this quote since it is the 2nd sequential
						// double-quote which means that it is escaped
						openQuote = -1
						continue
					} else {
						// close the quotes...
						openQuote = -1
					}
				} else {
					// this is an opening quote
					openQuote = i
				}
			}
			field.WriteRune(c)
		case DELIMITER:
			// if this char was not escaped and we are not after an opening quote
			if nextEscapedCharIdx != i && openQuote == -1 {
				// end field and save it to row
				row = append(row, parseType(field.String()))
				// start new field
				field = bytes.NewBufferString("")
				fieldIdx++
			} else {
				field.WriteRune(c)
			}
		case NEW_LINE:
			// if this char was not escaped and we are not after an opening quote
			if nextEscapedCharIdx != i && openQuote == -1 {
				// end field and save it to row
				row = append(row, parseType(field.String()))
				// start new field
				field = bytes.NewBufferString("")
				fieldIdx++
				// end row and save it to rows
				rows = append(rows, row)
				// start new row
				row = []any{}
				rowIdx++
			} else {
				field.WriteRune(c)
			}
		default:
			field.WriteRune(c)
		}
	}
	if field.Len() > 0 {
		row = append(row, parseType(field.String()))
		rows = append(rows, row)
	}
	return rows
}

func parseType(field string) any {
	// if there is a quote in the field then it is a string
	if strings.Contains(field, "\"") {
		return field
	} else if strings.ToUpper(field) == "TRUE" {
		return true
	} else if strings.ToUpper(field) == "FALSE" {
		return false
	} else {
		intVal, err := strconv.Atoi(field)
		// if the field can be interpreted as an integer without an error, then return an int
		if err == nil {
			return intVal
		}
		floatVal, err := strconv.ParseFloat(strings.TrimSpace(field), 64)
		// if the field can be interpreted as a float without an error, then return a float
		if err == nil {
			return floatVal
		}
	}
	return field
}
