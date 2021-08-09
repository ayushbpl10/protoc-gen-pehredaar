package main

import (
	"regexp"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

func reverseString(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseString(input[1:]), input[0])
}

func ReturnFields(field pgs.Field, fields []pgs.Field) []pgs.Field {

	if field.Type().IsEmbed() {
		fields = append(fields, ReturnFields(field, fields)...)
	}
	fields = append(fields, field)
	return fields
}
