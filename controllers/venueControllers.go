package controllers

import (
	"context"
	"fmt"
	constants "instix_auth/constants"
	"instix_auth/database"
	helper "instix_auth/helpers"
	models "instix_auth/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var venueCollection *mongo.Collection = database.OpenCollection(database.Client, constants.VENUEDATABASE)

func CreateVenue() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var venue models.Venue

		defer cancel()

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := c.Bind(&venue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Form binded successfully")

		validationErr := validate.Struct(venue)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		fmt.Println("Validation checked successfully")

		venue.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		venue.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		venue.ID = primitive.NewObjectID()
		venue.Creator = c.GetString("first_name")

		file, _ := c.FormFile("images")
		fp := constants.VenueImageDir + "/" + venue.ID.Hex()
		imageURL := helper.GetImageURL(file, venue.ID.Hex(), fp, c)

		if imageURL == constants.IMAGE_NOT_UPLOADED {
			c.JSON(http.StatusInternalServerError, gin.H{"error": imageURL})
			return
		}

		venue.ImageURL = imageURL

		resultInsertionNumber, insertErr := venueCollection.InsertOne(ctx, venue)
		if insertErr != nil {
			msg := constants.VENUE_NOT_CREATED
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}
