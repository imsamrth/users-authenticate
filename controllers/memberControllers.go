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
			log.Print("Error finding members", err)
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

		c.JSON(http.StatusOK, gin.H{"count": len(members), "members": members})
	}
}

func GetMember() gin.HandlerFunc {
	return func(c *gin.Context) {

		mid := c.Param("member_id")
		var member models.Member

		findOpts := options.FindOne()
		projection := bson.D{
			{Key: "por", Value: 1},
			{Key: "body", Value: 1},
			{Key: "session", Value: 1},
			{Key: "level", Value: 1},
			{Key: "uid", Value: 1},
		}
		findOpts.SetProjection(projection)

		_id, err := primitive.ObjectIDFromHex(mid)
		err = memberCollection.FindOne(c, bson.D{{Key: "_id", Value: _id}}, findOpts).Decode(&member)

		if err != nil {
			log.Print("Error finding member", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, member)
	}
}

func AddMember() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		bid := c.Param("body_id")
		var payload struct {
			Session  string
			Password string
			Member   models.Member `json:"member"`
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

		id, _ := primitive.ObjectIDFromHex(bid)
		var foundBody models.Body
		err := bodyCollection.FindOne(c, bson.M{"_id": id}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: 1}})).Decode(&foundBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		passwordIsValid, msg := helper.VerifyPassword(payload.Password, *foundBody.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		session := payload.Session

		defer cancel()

		member := payload.Member
		member.ID = primitive.NewObjectID()
		member.BID = bid
		member.Session = session

		resultInsertionNumber, insertErr := memberCollection.InsertOne(ctx, member)
		if insertErr != nil {
			msg := "Member not created"
			fmt.Println(insertErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func UpdateMember() gin.HandlerFunc {

	return func(c *gin.Context) {
		mid := c.Param("member_id")

		var payload struct {
			Bid      string        `json:"bid"`
			Password string        `json:"password"`
			Member   models.Member `json:"member"`
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

		if statusCode, err := helper.IsMemberOfBody(mid, payload.Bid, memberCollection); err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		_id, _ := primitive.ObjectIDFromHex(mid)
		filter := bson.D{{Key: "_id", Value: _id}}
		update := bson.D{{Key: "$set", Value: payload.Member}}

		result, err := memberCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := "Member has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No Member is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, "Member updated successfully")
	}
}

func DeleteMember() gin.HandlerFunc {

	return func(c *gin.Context) {
		mid := c.Param("member_id")

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

		if statusCode, err := helper.IsMemberOfBody(mid, payload.Bid, memberCollection); err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		_id, _ := primitive.ObjectIDFromHex(mid)

		opts := options.Delete().SetCollation(&options.Collation{})

		res, err := memberCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts)
		if err != nil || (res.DeletedCount == 0) {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": constants.MEMBER_NOT_DELETED})
			return
		}

		fmt.Println(res.DeletedCount)
		c.JSON(http.StatusAccepted, gin.H{"msg": "Deleted successfully"})
	}
}
