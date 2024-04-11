package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func formatDate(date time.Time) string {
	return date.Format("01/02/2006")
}

func formatSplit(split bool) string {
	out := "No"
	if split {
		out = "Yes"
	}
	return out
}

func formatCost(cost float32) string {
	return fmt.Sprintf("$%.2f", cost)
}

func capitalizeAndTrim(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.TrimSpace(string(unicode.ToUpper(rune(name[0]))) + name[1:])
}

func isValidName(name string) bool {
	match, err := regexp.Match(`^[a-zA-Z'-]+$`, []byte(name))
	if err != nil {
		return false
	}
	return match
}
