package controllers

import (
	"context"
	"fmt"
	"instix_auth/constants"
	"instix_auth/database"
	models "instix_auth/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
