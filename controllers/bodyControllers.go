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

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bodyCollection *mongo.Collection = database.OpenCollection(database.Client, constants.BODIESDATABASE)

func CreateBody() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var body models.Body
		if err := c.Bind(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(body)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		body.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		body.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		body.ID = primitive.NewObjectID()

		statusCode, err := helper.ValidateUsername(*&body.Username, body.ID.Hex(), profileCollection, ctx)
		if err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		password := helper.HashPassword(*body.Password)
		body.Password = &password

		count, err := bodyCollection.CountDocuments(ctx, bson.M{"name": body.Name})

		fmt.Print("Checking does name already exists")
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.BODY_NAME_ALREADY_TAKEN})
		}

		body.Verified = false

		file, _ := c.FormFile("images")
		fp := constants.BodyLogoDir + "/" + body.ID.Hex()
		sp := constants.BodyLogoURL + "/" + body.ID.Hex()
		imageURL := helper.GetImageURL(file, body.ID.Hex(), fp, sp, c)

		if imageURL == constants.IMAGE_NOT_UPLOADED {
			c.JSON(http.StatusInternalServerError, gin.H{"error": imageURL})
			return
		}

		body.ImageURL = imageURL

		resultInsertionNumber, insertErr := bodyCollection.InsertOne(ctx, body)
		if insertErr != nil {
			msg := fmt.Sprintf("Body was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}
