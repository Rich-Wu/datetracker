package main

import (
	"fmt"
	"time"
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
