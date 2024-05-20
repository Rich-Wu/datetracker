package main

const (
	APP_NAME       = "datetracker"
	DATES_TABLE    = "dates"
	DEV            = "DEV"
	DEV_MONGO      = "mongodb://localhost:27017"
	DEV_PORT       = 8080
	DEV_SECRET     = "secret"
	EMAILS_TABLE   = "emails"
	IMAGES_PATH    = "./images"
	PROD           = "PROD"
	SESSIONS_STORE = "sessions"
	SESSION_TTL    = 3600
	USERS_TABLE    = "users"
)

var (
	SUPPORTED_IMAGE_TYPES = make(map[string]bool)
)

func init() {
	SUPPORTED_IMAGE_TYPES["image/jpeg"] = true
	SUPPORTED_IMAGE_TYPES["image/png"] = true
}
