package controllers

import (
	"context"
	"fmt"
	"instix_auth/database"
	helper "instix_auth/helpers"
	models "instix_auth/models"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

var (
	domName = os.Getenv("DomainName")
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Hash is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		statusCode, err := helper.ValidateEmail(*user.Email)
		fmt.Println("Validation of email is called")
		if err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		fmt.Print("Email Validation is checked successfully")
		password := HashPassword(*user.Password)
		user.Password = &password

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		fmt.Print("Checking is email already registered")
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the user email"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email is already registered"})
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		user.Activate = false

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func ActivateGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
		emailVerRandRune := make([]rune, 64)
		// creat a random slice of runes (characters) to create our emailVerPassword (random string of characters)
		for i := 0; i < 64; i++ {
			emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
		}
		fmt.Println("emailVerRandRune:", emailVerRandRune)
		emailVerPassword := string(emailVerRandRune)
		fmt.Println("emailVerPassword:", emailVerPassword)
		var emailVerPWhash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		emailVerPWhash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("bcrypt err:", err)
			return
		}
		fmt.Println("emailVerPWhash:", emailVerPWhash)
		user.VerHash = string(emailVerPWhash)

		id := user.ID
		filter := bson.D{{"_id", id}}
		update := bson.D{{"$set", bson.D{{"ver_hash", user.VerHash}, {"activate", false}}}}

		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			panic(err)
		}
		subject := "Email verification for your InstiX account"
		// Add handler to detect the server

		HTMLbody :=
			`<html> 
				<h1>Click Link to Veify Email</h1>
				<a href="` + domName + `/activate/` + user.User_id + `/` + emailVerPassword + `">click to verify email</a>
			</html>`

		verificationLink := domName + `/activate/` + user.User_id + `/` + emailVerPassword
		err = helper.SendEmail(*user.First_name, *user.Email, subject, HTMLbody, verificationLink)

		if err != nil {
			fmt.Println("issue sending verification email")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func EmailverGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		verHash := c.Param("ver_hash")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		isverified, msg := VerifyPassword(verHash, user.VerHash)
		if isverified {
			id := user.ID
			filter := bson.D{{"_id", id}}
			update := bson.D{{"$set", bson.D{{"activate", true}}}}

			result, err := userCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.JSON(http.StatusOK, result)
		} else {
			fmt.Println(msg)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect hash code"})
		}
		defer cancel()
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
