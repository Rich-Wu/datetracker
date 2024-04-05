package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var db *mongo.Database
	var Stage string = os.Getenv("STAGE")
	var db_uri string
	var client *mongo.Client
	// var users *mongo.Collection
	var dates *mongo.Collection
	var router *gin.Engine

	if Stage == DEV {
		db_uri = DEV_MONGO
	}

	clientOptions := options.Client().ApplyURI(db_uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully connected to mongodb at %s\n", clientOptions.GetURI())

	db = client.Database(APP_NAME)
	// users = db.Collection(USERS_TABLE)
	dates = db.Collection(DATES_TABLE)

	router = gin.Default()
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/login.html")
	})
	router.GET("/date/new", func(c *gin.Context) {
		c.File("./static/date.html")
	})
	router.GET("/dates", func(c *gin.Context) {
		// TODO
	})
	api := router.Group("/api")
	{
		api.GET("/dates", func(c *gin.Context) {
			findOptions := options.Find()
			// Sort by the date of occurrence, descending and then recency of insertion for tiebreaking
			findOptions.SetSort(bson.D{{Key: "date", Value: -1}, {Key: "createdAt", Value: -1}})
			// TODO: allow changing limit in query params
			findOptions.SetLimit(20)

			cursor, err := dates.Find(context.Background(), bson.D{}, findOptions)
			if err != nil {
				log.Println("Error finding documents:", err)
				c.AbortWithError(http.StatusConflict, err)
			}
			defer cursor.Close(context.Background())

			foundDates := []Date{}
			for cursor.Next(context.Background()) {
				var result Date
				if err := cursor.Decode(&result); err != nil {
					fmt.Println("Error decoding document:", err)
					c.AbortWithError(http.StatusInternalServerError, err)
				}
				foundDates = append(foundDates, result)
			}

			c.JSON(http.StatusOK, foundDates)
		})

		api.POST("/date/new", func(c *gin.Context) {
			c.Request.ParseForm()
			formData := c.Request.Form

			cost, _ := strconv.ParseFloat(formData["cost"][0], 32)
			age, _ := strconv.ParseInt(formData["age"][0], 10, 32)
			date, _ := time.Parse(time.DateOnly, formData["date"][0])
			split, _ := strconv.ParseBool(formData["split"][0])

			newDate := &Date{
				FirstName:  formData["first_name"][0],
				LastName:   formData["last_name"][0],
				Ethnicity:  formData["ethnicity"][0],
				Occupation: formData["occupation"][0],
				Place:      formData["place"][0],
				TypeOfDate: formData["type_of_date"][0],
				Cost:       float32(cost),
				Result:     formData["how_ended"][0],
				Age:        int32(age),
				Date:       date,
				Split:      split,
				CreatedAt:  time.Now(),
			}

			_, err := dates.InsertOne(context.Background(), newDate)
			if err != nil {
				log.Fatalln("insertion to db failed", err)
			}

			c.Redirect(http.StatusFound, "/dates")
		})
	}

	router.Run(":8080")
}
