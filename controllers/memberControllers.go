package controllers

import (
	"context"
	"fmt"
	"instix_auth/constants"
	"instix_auth/database"
	helper "instix_auth/helpers"
	models "instix_auth/models"
	"log"
	"net/http"
	"time"

	"github.com/ajclopez/mgs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var memberCollection *mongo.Collection = database.OpenCollection(database.Client, constants.MEMEBERDATABASE)

func CreateCouncil() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		bid := c.Param("body_id")
		var payload struct {
			Session  string
			Password string
			Members  []models.Member `json:"members"`
		}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := primitive.ObjectIDFromHex(bid)
		var foundBody models.Body
		err = bodyCollection.FindOne(c, bson.M{"_id": id}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: 1}})).Decode(&foundBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		passwordIsValid, msg := helper.VerifyPassword(payload.Password, *foundBody.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		validationErr := validate.Struct(payload)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		session := payload.Session
		err = helper.IsAlreadyCouncilCreated(session, bid, memberCollection)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cancel()

		var members []interface{}
		for _, member := range payload.Members {
			member.ID = primitive.NewObjectID()
			member.BID = bid
			member.Session = session
			members = append(members, member)
		}

		resultInsertionNumber, insertErr := memberCollection.InsertMany(ctx, members)
		if insertErr != nil {
			msg := fmt.Sprint("Council was not created")
			fmt.Println(members...)
			fmt.Println(insertErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func GetMembers() gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Request.URL.RawQuery

		var members []models.Member

		opts := mgs.FindOption()
		// Set max limit to restrict the number of results returned
		opts.SetMaxLimit(100)
		result, err := mgs.MongoGoSearch(query, opts)
		if err != nil {
			//invalid query
			log.Print("Invalid query", err)
			c.JSON(http.StatusInternalServerError, members)
			return
		}

		findOpts := options.Find()
		findOpts.SetLimit(result.Limit)
		findOpts.SetSkip(result.Skip)
		findOpts.SetSort(result.Sort)

		projection := bson.D{
			{Key: "por", Value: 1},
			{Key: "body", Value: 1},
			{Key: "session", Value: 1},
			{Key: "level", Value: 1},
			{Key: "uid", Value: 1},
		}
		findOpts.SetProjection(projection)

		cur, err := memberCollection.Find(c, result.Filter, findOpts)
		if err != nil {
			log.Print("Error finding products", err)
			c.JSON(http.StatusInternalServerError, members)
			return
		}
		for cur.Next(c) {
			var member models.Member
			err := cur.Decode(&member)
			if err != nil {
				c.JSON(http.StatusInternalServerError, members)
				return
			}
			members = append(members, member)
		}

		c.JSON(http.StatusOK, gin.H{"count": len(members), "products": members})
	}
}
