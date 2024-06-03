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

var eventCollection *mongo.Collection = database.OpenCollection(database.Client, constants.EVENTDATABASE)

func CreateEvent() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		bid := c.Param("body_id")
		var payload struct {
			Session  string
			Password string
			Event    models.Event `json:"events"`
		}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(payload)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
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

		defer cancel()

		resultInsertionNumber, insertErr := eventCollection.InsertOne(ctx, payload.Event)
		if insertErr != nil {
			msg := fmt.Sprint("Event was not created")
			fmt.Println(insertErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Request.URL.RawQuery

		var events []models.Event

		opts := mgs.FindOption()
		// Set max limit to restrict the number of results returned
		opts.SetMaxLimit(20)
		result, err := mgs.MongoGoSearch(query, opts)
		if err != nil {
			//invalid query
			log.Print("Invalid query", err)
			c.JSON(http.StatusInternalServerError, events)
			return
		}

		findOpts := options.Find()
		findOpts.SetLimit(result.Limit)
		findOpts.SetSkip(result.Skip)
		findOpts.SetSort(result.Sort)

		cur, err := eventCollection.Find(c, result.Filter, findOpts)
		if err != nil {
			log.Print("Error finding events", err)
			c.JSON(http.StatusInternalServerError, events)
			return
		}
		for cur.Next(c) {
			var event models.Event
			err := cur.Decode(&event)
			if err != nil {
				c.JSON(http.StatusInternalServerError, events)
				return
			}
			events = append(events, event)
		}

		c.JSON(http.StatusOK, gin.H{"count": len(events), "events": events})
	}
}

func GetEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

		eid := c.Param("event_id")
		var event models.Event

		findOpts := options.FindOne()

		_id, err := primitive.ObjectIDFromHex(eid)
		err = eventCollection.FindOne(c, bson.D{{Key: "_id", Value: _id}}, findOpts).Decode(&event)

		if err != nil {
			log.Print("Error finding event", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, event)
	}
}

func UpdateEvent() gin.HandlerFunc {

	return func(c *gin.Context) {
		eid := c.Param("event_id")

		var payload struct {
			Bid      string       `json:"bid"`
			Password string       `json:"password"`
			Event    models.Event `json:"event"`
		}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(payload)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if err := helper.VerifyBodyId(payload.Bid, payload.Password, bodyCollection); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if statusCode, err := helper.IsEventOfBody(eid, payload.Bid, eventCollection); err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		_id, _ := primitive.ObjectIDFromHex(eid)
		filter := bson.D{{Key: "_id", Value: _id}}
		update := bson.D{{Key: "$set", Value: payload.Event}}

		result, err := eventCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := "event has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No event is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, "event updated successfully")
	}
}

func DeleteEvent() gin.HandlerFunc {

	return func(c *gin.Context) {
		eid := c.Param("event_id")

		var payload struct {
			Bid      string `json:"bid"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(payload)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if err := helper.VerifyBodyId(payload.Bid, payload.Password, bodyCollection); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if statusCode, err := helper.IsEventOfBody(eid, payload.Bid, eventCollection); err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		_id, _ := primitive.ObjectIDFromHex(eid)

		opts := options.Delete().SetCollation(&options.Collation{})

		res, err := eventCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts)
		if err != nil || (res.DeletedCount == 0) {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": constants.EVENT_NOT_DELETED})
			return
		}

		fmt.Println(res.DeletedCount)
		c.JSON(http.StatusAccepted, gin.H{"msg": "Deleted successfully"})
	}
}
