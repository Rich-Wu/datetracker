package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var db *mongo.Database
	var Stage string = os.Getenv("STAGE")
	var db_uri string
	var secret string
	var client *mongo.Client
	var usersCollection *mongo.Collection
	var datesCollection *mongo.Collection
	var sessionsCollection *mongo.Collection
	var router *gin.Engine

	if Stage == DEV {
		db_uri = DEV_MONGO
		secret = DEV_SECRET
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
	usersCollection = db.Collection(USERS_TABLE)
	datesCollection = db.Collection(DATES_TABLE)
	sessionsCollection = db.Collection(SESSIONS_STORE)

	// Sets username to be unique
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // Create a unique index on the "username" field
		Options: options.Index().SetUnique(true),
	}

	// Create the unique index
	_, err = usersCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	router = gin.Default()
	sessionsStore := mongodriver.NewStore(sessionsCollection, 3600, true, []byte(secret))
	router.Use(sessions.Sessions("session", sessionsStore))
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get("user")
		username := session.Get("username")
		count := session.Get("count")
		if count == nil {
			count = 0
		}
		log.Println("userId:", userId)
		log.Println("username:", username)
		log.Println("count:", count)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"username": username,
			"count":    count,
			"session":  session,
		})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{})
	})
	router.GET("/date/new", func(c *gin.Context) {
		session := sessions.Default(c)
		log.Println(session.Get("username"))
		c.HTML(http.StatusOK, "date.tmpl", gin.H{
			"username": session.Get("username"),
		})
	})
	router.GET("/dates", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		user := session.Get("user")
		test := "teststring"

		c.HTML(http.StatusOK, "dates.tmpl", gin.H{
			"test":     test,
			"username": username,
			"user":     user,
		})
	})
	router.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		c.Request.ParseForm()
		filter := bson.D{{Key: "username", Value: c.PostForm("username")}}
		userResult := usersCollection.FindOne(context.Background(), filter, options.FindOne())

		foundUser := &User{}
		userResult.Decode(foundUser)

		err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(c.PostForm("password")))
		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
		}

		session.Set("user", foundUser.ID.Hex())
		session.Set("username", foundUser.UserName)
		session.Save()

		c.Redirect(303, "/")
	})
	api := router.Group("/api")
	{
		api.GET("/dates", func(c *gin.Context) {
			findOptions := options.Find()
			// Sort by the date of occurrence, descending and then recency of insertion for tiebreaking
			findOptions.SetSort(bson.D{{Key: "date", Value: -1}, {Key: "createdAt", Value: -1}})
			// TODO: allow changing limit in query params
			findOptions.SetLimit(20)

			cursor, err := datesCollection.Find(context.Background(), bson.D{}, findOptions)
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

		api.POST("/user/new", func(c *gin.Context) {
			c.Request.ParseForm()
			session := sessions.Default(c)
			password := c.PostForm("password")

			// Hash the password before storing it
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			}

			newUser := &User{
				UserName:  c.PostForm("username"),
				FirstName: c.PostForm("first_name"),
				LastName:  c.PostForm("last_name"),
				Password:  string(hashedPassword),
			}

			// TODO: Do something with the user for session storage
			_, err = usersCollection.InsertOne(context.Background(), newUser)
			if err != nil {
				log.Println("insertion to db failed", err)
				c.JSON(http.StatusUnprocessableEntity, errors.New(err.Error()))
			}
			session.Set("username", c.PostForm("username"))
			session.Save()
			log.Println(session)
			log.Println(session.Get("username"))

			c.Redirect(http.StatusSeeOther, "/date/new")
		})

		api.POST("/date/new", func(c *gin.Context) {
			c.Request.ParseForm()
			session := sessions.Default(c)

			cost, _ := strconv.ParseFloat(c.PostForm("cost"), 32)
			age, _ := strconv.ParseInt(c.PostForm("age"), 10, 32)
			date, _ := time.Parse(time.DateOnly, c.PostForm("date"))
			split, _ := strconv.ParseBool(c.PostForm("split"))
			owner := session.Get("user")
			if owner == nil {
				c.AbortWithStatus(http.StatusForbidden)
			}

			newDate := &Date{
				Owner:      owner.(primitive.ObjectID),
				FirstName:  c.PostForm("first_name"),
				LastName:   c.PostForm("last_name"),
				Ethnicity:  c.PostForm("ethnicity"),
				Occupation: c.PostForm("occupation"),
				Place:      c.PostForm("place"),
				TypeOfDate: c.PostForm("type_of_date"),
				Cost:       float32(cost),
				Result:     c.PostForm("result"),
				Age:        int32(age),
				Date:       date,
				Split:      split,
				CreatedAt:  time.Now(),
			}

			_, err := datesCollection.InsertOne(context.Background(), newDate)
			if err != nil {
				log.Fatalln("insertion to db failed", err)
			}

			c.Redirect(http.StatusFound, "/dates")
		})
	}

	router.Run(":8080")
}
