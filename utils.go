package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func dateString(date time.Time) string {
	year, month, day := date.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

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

func isValidName(name string) bool {
	match, err := regexp.Match(`^[a-zA-Z'-]+$`, []byte(name))
	if err != nil {
		return false
	}
	return match
}

func sanitizeUsername(username string) string {
	// Trim whitespace
	username = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(username, "")
	// Remove special characters and keep Chinese characters
	re := regexp.MustCompile(`[^\w\s\p{Han}]`)
	username = re.ReplaceAllString(username, "")

	return username
}

func toSingleString(stringArr []string) string {
	return strings.Join(stringArr, ",") + ","
}

func getFiletypeFromMime(mime string) string {
	parts := strings.Split(mime, "/")
	return parts[len(parts)-1]
}

func getHex(objectId primitive.ObjectID) string {
	return objectId.Hex()
}

// Use this function as shorthand for c.HTML(code, "error.tmpl", *custom error data)
func renderError(c *gin.Context, code int) {
	var errorMap = make(map[int]*gin.H)
	errorMap[400] = &gin.H{
		"errorCode": http.StatusBadRequest,
		"errorName": "Bad Request",
		"errorMessage": "Oops! It seems like something went wrong with your request. Our server couldn't " +
			"understand what you're asking for. This might be due to a mistake in the way the request was " +
			"formatted or invalid data being sent.",
	}
	errorMap[403] = &gin.H{
		"errorCode": http.StatusForbidden,
		"errorName": "Forbidden Page",
		"errorMessage": "Uh oh! You seem to have stumbled somewhere you shouldn't have. Please make sure you have permission " +
			"to access the requested resource. If you are sure you should be able to access this page, please contact us for " +
			"assistance.",
	}
	errorMap[404] = &gin.H{
		"errorCode": http.StatusNotFound,
		"errorName": "Not Found",
		"errorMessage": "Oh no! We couldn't find the resource you were looking for. Please try searching for your resource " +
			"again, or if you are sure the resource exists, please report this issue to us and we'll find out what's wrong. ",
	}
	errorMap[409] = &gin.H{
		"errorCode": http.StatusConflict,
		"errorName": "Resource Conflict",
		"errorMessage": "There's a conflict a-brewing. A requested resource was being accessed and cannot be " +
			"accessed at this time. Please retry your request at a later time.",
	}
	errorMap[422] = &gin.H{
		"errorCode": http.StatusUnprocessableEntity,
		"errorName": "Unprocessable Content",
		"errorMessage": "Uh oh. Something went wrong during the last requested operation. Please contact us with details " +
			"about this error and we'll do our best to figure out what went wrong",
	}
	errorMap[500] = &gin.H{
		"errorCode": http.StatusInternalServerError,
		"errorName": "Internal Server Error",
		"errorMessage": "Oh no! It seems like something went wrong on our end. This is due either to something unexpected " +
			"happening or may be a temporary issue. Please retry this request at a later time or you can contact us " +
			"with details of what happened.",
	}
	c.HTML(code, "error.tmpl", *errorMap[code])
}
