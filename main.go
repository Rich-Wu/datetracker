package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/appengine/v2/mail"
)

func main() {
	var db *mongo.Database
	var Stage string = os.Getenv("STAGE")
	var port int32
	var db_uri string
	var secret string
	var client *mongo.Client
	var emailsCollection *mongo.Collection
	var usersCollection *mongo.Collection
	var datesCollection *mongo.Collection
	var sessionsCollection *mongo.Collection
	var router *gin.Engine
	// String processor
	var caser cases.Caser

	if Stage == DEV {
		db_uri = DEV_MONGO
		secret = DEV_SECRET
		port = DEV_PORT
	} else if Stage == PROD {
		db_uri = os.Getenv("DB_URI")
		secret = os.Getenv("SECRET")
		portNum, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
		port = int32(portNum)
	} else {
		log.Fatalln("stage was not set correctly. Check configurations and environments.")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(db_uri).SetServerAPIOptions(serverAPI)

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
	emailsCollection = db.Collection(EMAILS_TABLE)
	usersCollection = db.Collection(USERS_TABLE)
	datesCollection = db.Collection(DATES_TABLE)
	sessionsCollection = db.Collection(SESSIONS_STORE)

	// Sets username to be unique
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // Create a unique index on the "username" field
		Options: options.Index().SetUnique(true),
	}
	// Create the unique index
	name, err := usersCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name of Index Created: " + name)
	// Create unique email index for emails collection
	indexModel = mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	name, err = emailsCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name of Index Created: " + name)

	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "ownerId", Value: 1}},
	}
	name, err = datesCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name of Index Created: " + name)

	caser = cases.Title(language.AmericanEnglish)

	router = gin.Default()
	sessionsStore := mongodriver.NewStore(sessionsCollection, 3600, true, []byte(secret))
	router.Use(sessions.Sessions("session", sessionsStore))
	router.SetFuncMap(template.FuncMap{
		"dateString":     dateString,
		"formatCost":     formatCost,
		"formatDate":     formatDate,
		"formatSplit":    formatSplit,
		"getHex":         getHex,
		"toSingleString": toSingleString,
	})
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*.tmpl")

	// router.GET("/test", func(c *gin.Context) {

	// })
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/termsOfService", func(c *gin.Context) {
		c.File("./static/termsOfService.txt")
	})
	router.GET("/privacyPolicy", func(c *gin.Context) {
		c.File("./static/privacyPolicy.txt")
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	router.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		c.Redirect(http.StatusFound, "/")
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{})
	})
	router.GET("/date/new", func(c *gin.Context) {
		session := sessions.Default(c)
		c.HTML(http.StatusOK, "date.tmpl", gin.H{
			"username": session.Get("username"),
		})
	})
	router.GET("/dates", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		user := session.Get("user")
		if user == nil || username == nil {
			renderError(c, http.StatusForbidden)
			return
		}
		findOptions := options.Find()
		// Sort by date of occurrence, descending, then by time of entry for tiebreaking
		findOptions.SetSort(bson.D{{Key: "date", Value: -1}, {Key: "createdAt", Value: -1}})
		// TODO: allow changing limit in query params
		findOptions.SetLimit(50)

		userId, err := primitive.ObjectIDFromHex(user.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("could not parse user"))
		}

		cursor, err := datesCollection.Find(context.Background(), bson.D{{Key: "ownerId", Value: userId}}, findOptions)
		if err != nil {
			log.Println("Error finding focuments:", err)
			c.AbortWithError(http.StatusConflict, err)
		}
		defer cursor.Close(context.Background())

		dates := []Date{}
		for cursor.Next(context.Background()) {
			var result Date
			if err := cursor.Decode(&result); err != nil {
				log.Println("Error decoding document:", err)
				renderError(c, http.StatusInternalServerError)
			}
			dates = append(dates, result)
		}
		c.HTML(http.StatusOK, "dates.tmpl", gin.H{
			"username": username,
			"user":     user,
			"dates":    dates,
		})
	})
	router.GET("/dates/:username", func(c *gin.Context) {
		username := c.Param("username")
		userResult := usersCollection.FindOne(context.Background(), bson.D{{Key: "username", Value: username}}, options.FindOne())
		if userResult.Err() != nil {
			c.Redirect(http.StatusSeeOther, "/dates")
		}

		foundUser := &User{}
		userResult.Decode(foundUser)

		findOptions := options.Find()
		// Sort by date of occurrence, descending, then by time of entry for tiebreaking
		findOptions.SetSort(bson.D{{Key: "date", Value: -1}, {Key: "createdAt", Value: -1}})
		// TODO: allow changing limit in query params
		findOptions.SetLimit(50)

		cursor, err := datesCollection.Find(context.Background(), bson.D{{Key: "ownerId", Value: foundUser.ID}}, findOptions)
		if err != nil {
			log.Println("Error finding focuments:", err)
			c.AbortWithError(http.StatusConflict, err)
		}
		defer cursor.Close(context.Background())

		dates := []Date{}
		for cursor.Next(context.Background()) {
			var result Date
			if err := cursor.Decode(&result); err != nil {
				fmt.Println("Error decoding document:", err)
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			dates = append(dates, result)
		}
		c.HTML(200, "dates.tmpl", gin.H{
			"username": username,
			"dates":    dates,
		})
	})
	router.GET("/date/edit/:id", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		user := session.Get("user")
		if user == nil || username == nil {
			renderError(c, http.StatusForbidden)
			return
		}
		idHex := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			renderError(c, http.StatusNotFound)
			return
		}
		dateResult := datesCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}, options.FindOne())
		if dateResult.Err() != nil {
			log.Println("Error when finding a date result:", err)
			renderError(c, http.StatusNotFound)
			return
		}

		date := &Date{}
		err = dateResult.Decode(date)
		if err != nil {
			log.Println("Decoding date result failed:", err)
			renderError(c, http.StatusInternalServerError)
			return
		}
		if date.OwnerId.Hex() != user {
			log.Println("Mismatched user and date")
			renderError(c, http.StatusForbidden)
			return
		}

		c.HTML(200, "date.tmpl", gin.H{
			"id":   idHex,
			"date": date,
		})
	})
	router.POST("/email", func(c *gin.Context) {
		c.Request.ParseForm()

		email := &Email{}
		if err := c.ShouldBindWith(email, binding.Form); err != nil {
			log.Println("An invalid email was provided:", email.Address, err, email.Name)
			renderError(c, http.StatusBadRequest)
			return
		}
		email.SignupTime = time.Now().UTC()

		id, err := emailsCollection.InsertOne(context.Background(), email)
		if err != nil {
			log.Println("Error occurred when saving an email to db:", err)
			renderError(c, http.StatusUnprocessableEntity)
			return
		}

		var confirmMail bytes.Buffer
		tmpl, _ := template.ParseGlob("./templates/email/confirm.tmpl")
		err = tmpl.ExecuteTemplate(&confirmMail, "email/confirm.tmpl", gin.H{
			"email": email,
			"id":    id.InsertedID.(primitive.ObjectID).Hex(),
		})
		if err != nil {
			log.Println("An error occurred when rendering:", err)
			return
		}
		msg := &mail.Message{
			Sender:   "richie1988@gmail.com",
			To:       []string{fmt.Sprintf("%s <%s>", email.Name, email.Address)},
			Subject:  "Thank you for your interest in littleblackbook",
			HTMLBody: confirmMail.String(),
		}

		if err := mail.Send(c, msg); err != nil {
			log.Println("An error occurred sending a response mail:", err)
		}

		c.HTML(http.StatusOK, "confirm.tmpl", gin.H{
			"email": email,
		})
	})
	router.GET("/email/unsubscribe/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			log.Println("Bad visitor to unsubscribe route")
			renderError(c, http.StatusBadRequest)
			return
		}

		deleteResult, err := emailsCollection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}}, options.Delete())
		if err != nil {
			log.Println("Error occurred when deleting an email:", err)
			renderError(c, http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, deleteResult)
	})
	router.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		c.Request.ParseForm()
		filter := bson.D{{Key: "username", Value: c.PostForm("username")}}
		userResult := usersCollection.FindOne(context.Background(), filter, options.FindOne())
		if userResult.Err() != nil {
			renderError(c, http.StatusForbidden)
		}

		foundUser := &User{}
		err := userResult.Decode(foundUser)
		if err != nil {
			log.Printf("An error occurred decoding a userResult")
			renderError(c, http.StatusInternalServerError)
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(c.PostForm("password")))
		if err != nil {
			renderError(c, http.StatusForbidden)
		}

		session.Set("user", foundUser.ID.Hex())
		session.Set("username", foundUser.UserName)
		session.Save()

		c.Redirect(http.StatusFound, "/dates")
	})
	api := router.Group("/api")
	{
		api.POST("/date/:id", func(c *gin.Context) {
			session := sessions.Default(c)
			user := session.Get("user")
			c.Request.ParseForm()
			userId, err := primitive.ObjectIDFromHex(user.(string))
			if err != nil {
				log.Println("Invalid user in session", err)
				renderError(c, http.StatusNotFound)
				return
			}
			dateId := c.Param("id")
			dateObjId, err := primitive.ObjectIDFromHex(dateId)
			if err != nil {
				log.Println("Invalid date specified", err)
				renderError(c, http.StatusNotFound)
				return
			}
			date := &Date{}
			dateResult := datesCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: dateObjId}, {Key: "ownerId", Value: userId}}, options.FindOne())
			err = dateResult.Decode(date)
			if err != nil {
				log.Println("Error occurred when decoding result", err)
				renderError(c, http.StatusInternalServerError)
				return
			}

			var firstName string
			var lastName string
			age, _ := strconv.ParseInt(c.PostForm("age"), 10, 32)
			dateString, _ := time.Parse(time.DateOnly, c.PostForm("date"))

			if isValidName(c.PostForm("first_name")) {
				firstName = caser.String(c.PostForm("first_name"))
				lastName = caser.String(c.PostForm("last_name"))
			} else {
				renderError(c, http.StatusBadRequest)
				return
			}

			var runningTotal float32 = 0
			length := len(c.PostFormArray("place"))
			places := make([]*Place, length)
			for i := 0; i < length; i++ {
				cost, _ := strconv.ParseFloat(c.PostFormArray("cost")[i], 32)
				place := &Place{
					Place:       caser.String(c.PostFormArray("place")[i]),
					TypeOfPlace: c.PostFormArray("type_of_place")[i],
					Cost:        float32(cost),
				}
				places[i] = place
				runningTotal += float32(cost)
			}

			length = len(c.PostForm("ethnicity"))
			ethnicities := c.PostForm("ethnicity")[:length-1]
			ethnicitiesList := strings.Split(ethnicities, ",")
			split, _ := strconv.ParseBool(c.PostForm("split"))

			updatedDate := &Date{
				OwnerId:    userId,
				FirstName:  firstName,
				LastName:   lastName,
				Ethnicity:  ethnicitiesList,
				Occupation: caser.String(c.PostForm("occupation")),
				Places:     places,
				Split:      split,
				Cost:       runningTotal,
				Result:     c.PostForm("result"),
				Age:        int32(age),
				Date:       dateString,
				CreatedAt:  time.Now(),
			}

			_, err = datesCollection.ReplaceOne(context.Background(), bson.D{{Key: "_id", Value: date.ID}}, updatedDate, &options.ReplaceOptions{})
			if err != nil {
				log.Println("An error occurred when updating a date record:", err)
				renderError(c, http.StatusInternalServerError)
				return
			}

			c.Redirect(303, "/dates")
		})
		api.GET("/dates", func(c *gin.Context) {
			findOptions := options.Find()
			// Sort by the date of occurrence, descending and then recency of insertion for tiebreaking
			findOptions.SetSort(bson.D{{Key: "date", Value: -1}, {Key: "createdAt", Value: -1}})
			// TODO: allow changing limit in query params
			findOptions.SetLimit(50)

			cursor, err := datesCollection.Find(context.Background(), bson.D{}, findOptions)
			if err != nil {
				log.Println("Error finding documents:", err)
				renderError(c, http.StatusConflict)
				return
			}
			defer cursor.Close(context.Background())

			foundDates := []Date{}
			for cursor.Next(context.Background()) {
				var result Date
				if err := cursor.Decode(&result); err != nil {
					fmt.Println("Error decoding document:", err)
					renderError(c, http.StatusInternalServerError)
					return
				}
				foundDates = append(foundDates, result)
			}

			c.JSON(http.StatusOK, foundDates)
		})
		api.GET("/user/:id", func(c *gin.Context) {
			userId, err := primitive.ObjectIDFromHex(c.Param("id"))
			if err != nil {
				log.Println("An error occurred while looking up the user:", err)
				renderError(c, http.StatusBadRequest)
				return
			}
			userResult := usersCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: userId}}, options.FindOne())
			foundUser := &User{}
			userResult.Decode(&foundUser)
			c.JSON(200, foundUser)
		})
		api.POST("/user/new", func(c *gin.Context) {
			c.Request.ParseForm()
			session := sessions.Default(c)
			password := c.PostForm("password")
			// Hash the password before storing it
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				renderError(c, http.StatusInternalServerError)
				return
			}
			username := sanitizeUsername(c.PostForm("username"))

			// Profile picture processing
			var fileName string
			profilePic, err := c.FormFile("profilePic")
			if err != nil {
				if errors.Is(err, http.ErrMissingFile) {
					fileName = "default-pfp.jpeg"
				} else {
					log.Println("There was an unexpected error:", err)
					renderError(c, http.StatusBadRequest)
					return
				}
			} else {
				filetype := strings.ToLower(profilePic.Header.Get("Content-Type"))
				if !SUPPORTED_IMAGE_TYPES[strings.ToLower(profilePic.Header.Get("Content-Type"))] {
					log.Println("Unsupported filetype detected:", profilePic.Header.Get("Content-Type"))
					renderError(c, http.StatusBadRequest)
					return
				}
				// Save file to disk
				fileName = strings.ToLower(username) + "-pfp." + getFiletypeFromMime(filetype)
				f, err := os.OpenFile(filepath.Join(IMAGES_PATH, fileName), os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Println("Error creating file:", err)
					renderError(c, http.StatusInternalServerError)
					return
				}
				defer f.Close()
				file, err := profilePic.Open()
				if err != nil {
					log.Println("Error opening uploaded image:", err)
					renderError(c, http.StatusBadRequest)
					return
				}
				_, err = io.Copy(f, file)
				if err != nil {
					log.Println("Error copying uploaded image:", err)
					renderError(c, http.StatusBadRequest)
					return
				}
			}

			newUser := &User{
				UserName:   sanitizeUsername(c.PostForm("username")),
				FirstName:  caser.String(c.PostForm("first_name")),
				LastName:   caser.String(c.PostForm("last_name")),
				Password:   string(hashedPassword),
				ProfilePic: fileName,
			}

			result, err := usersCollection.InsertOne(context.Background(), newUser)
			if err != nil {
				renderError(c, http.StatusUnprocessableEntity)
				return
			}

			session.Set("username", c.PostForm("username"))
			session.Set("user", result.InsertedID.(primitive.ObjectID).Hex())
			session.Save()

			c.Redirect(http.StatusSeeOther, "/dates")
		})

		api.POST("/date/new", func(c *gin.Context) {
			var firstName string
			var lastName string
			c.Request.ParseForm()
			session := sessions.Default(c)

			age, _ := strconv.ParseInt(c.PostForm("age"), 10, 32)
			date, _ := time.Parse(time.DateOnly, c.PostForm("date"))

			owner := session.Get("user")
			if owner == nil {
				renderError(c, http.StatusForbidden)
				return
			}
			objId, err := primitive.ObjectIDFromHex(owner.(string))
			if err != nil {
				renderError(c, http.StatusForbidden)
				return
			}
			if isValidName(c.PostForm("first_name")) && isValidName(c.PostForm("last_name")) {
				firstName = caser.String(c.PostForm("first_name"))
				lastName = caser.String(c.PostForm("last_name"))
			} else {
				renderError(c, http.StatusBadRequest)
				return
			}

			var runningTotal float32 = 0
			length := len(c.PostFormArray("place"))
			places := make([]*Place, length)
			for i := 0; i < length; i++ {
				cost, _ := strconv.ParseFloat(c.PostFormArray("cost")[i], 32)
				place := &Place{
					Place:       caser.String(c.PostFormArray("place")[i]),
					TypeOfPlace: c.PostFormArray("type_of_place")[i],
					Cost:        float32(cost),
				}
				places[i] = place
				runningTotal += float32(cost)
			}

			length = len(c.PostForm("ethnicity"))
			ethnicities := c.PostForm("ethnicity")[:length-1]
			ethnicitiesList := strings.Split(ethnicities, ",")
			split, _ := strconv.ParseBool(c.PostForm("split"))

			newDate := &Date{
				OwnerId:    objId,
				FirstName:  firstName,
				LastName:   lastName,
				Ethnicity:  ethnicitiesList,
				Occupation: caser.String(c.PostForm("occupation")),
				Places:     places,
				Split:      split,
				Cost:       runningTotal,
				Result:     c.PostForm("result"),
				Age:        int32(age),
				Date:       date,
				CreatedAt:  time.Now(),
			}

			_, err = datesCollection.InsertOne(context.Background(), newDate)
			if err != nil {
				renderError(c, http.StatusInternalServerError)
				return
			}

			c.Redirect(http.StatusFound, "/dates")
		})
	}

	router.Run(":" + fmt.Sprint(port))
}
