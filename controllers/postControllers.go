package controllers

import (
	"context"
	"fmt"
	"instix_auth/constants"
	"instix_auth/database"
	models "instix_auth/models"
	"net/http"
	"time"

	"github.com/ajclopez/mgs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, constants.POSTDATABASE)

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var post models.Post

		defer cancel()
		if err := c.Bind(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(post)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		post.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		post.User_id = c.GetString("uid")
		post.ID = primitive.NewObjectID()
		post.Votes = []string{}
		post.Edited = true

		resultInsertionNumber, insertErr := postCollection.InsertOne(ctx, post)
		if insertErr != nil {
			fmt.Println(constants.POST_NOT_INSERTED)
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("post_id")
		uid := c.GetString("uid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var post models.Post

		id, err := primitive.ObjectIDFromHex(pid)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		opts := options.FindOne()
		projection := bson.D{
			{Key: "uid", Value: 1},
			{Key: "body", Value: 1},
			{Key: "link", Value: 1},
			{Key: "voted", Value: bson.D{{Key: "$in", Value: bson.A{uid, "$votes"}}}},
			{Key: "vote_count", Value: bson.D{{Key: "$size", Value: "$votes"}}},
			{Key: "tags", Value: 1},
			{Key: "edited", Value: 1},
		}
		opts.SetProjection(projection)
		err = postCollection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&post)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

func GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.RawQuery
		uid := c.GetString("uid")
		var posts []models.Post

		opts := mgs.FindOption()
		opts.SetMaxLimit(15)
		result, err := mgs.MongoGoSearch(query, opts)
		if err != nil {
			//invalid query
			fmt.Println("Invalid query", err)
			c.JSON(http.StatusInternalServerError, posts)
			return
		}

		findOpts := options.Find()
		findOpts.SetLimit(result.Limit)
		findOpts.SetSkip(result.Skip)
		findOpts.SetSort(result.Sort)

		projection := bson.D{
			{Key: "uid", Value: 1},
			{Key: "body", Value: 1},
			{Key: "link", Value: 1},
			{Key: "voted", Value: bson.D{{Key: "$in", Value: bson.A{uid, "$votes"}}}},
			{Key: "vote_count", Value: bson.D{{Key: "$size", Value: "$votes"}}},
			{Key: "tags", Value: 1},
			{Key: "edited", Value: 1},
		}
		findOpts.SetProjection(projection)

		cur, err := postCollection.Find(c, result.Filter, findOpts)
		if err != nil {
			fmt.Print("Error finding posts", err)
			c.JSON(http.StatusInternalServerError, posts)
			return
		}
		for cur.Next(c) {
			var post models.Post
			err := cur.Decode(&post)
			if err != nil {
				c.JSON(http.StatusInternalServerError, posts)
				return
			}
			posts = append(posts, post)
		}

		c.JSON(http.StatusOK, gin.H{"count": len(posts), "posts": posts})

	}
}

func ToggleVote() gin.HandlerFunc {

	return func(c *gin.Context) {
		pid := c.Param("post_id")
		uid := c.GetString("uid")

		_id, err := primitive.ObjectIDFromHex(pid) // convert params to //mongodb Hex ID
		if err != nil {
			fmt.Println(err.Error())
		}

		//TODO: Handle changes in image ( remove or add)
		filter := bson.D{{Key: "_id", Value: _id}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "votes", Value: bson.D{
					{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{
							{Key: "$in", Value: bson.A{uid, "$votes"}},
						}},
						{Key: "then", Value: bson.D{
							{Key: "$setDifference", Value: bson.A{
								"$votes", bson.A{uid},
							}},
						}},
						{Key: "else", Value: bson.D{
							{Key: "$concatArrays", Value: bson.A{"$votes", bson.A{uid}}},
						}},
					}},
				}},
			}},
		}

		result, err := postCollection.UpdateOne(context.TODO(), filter, mongo.Pipeline{update})

		if err != nil {
			fmt.Println(err.Error())
			msg := "Post has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No documents is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		fmt.Println(result)
		c.JSON(http.StatusOK, "Post updated successfully")
	}
}
