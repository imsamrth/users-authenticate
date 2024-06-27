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

	"github.com/ajclopez/mgs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		file, _ := c.FormFile("image")
		fp := constants.VenueImageDir + "/" + venue.ID.Hex()
		sp := constants.VenueImageURL + "/" + venue.ID.Hex()
		imageURL := helper.GetImageURL(file, venue.ID.Hex(), fp, sp, c)

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

		c.JSON(http.StatusOK, gin.H{"res": resultInsertionNumber, "id": venue.ID})
	}
}

func UpdateVenue() gin.HandlerFunc {

	return func(c *gin.Context) {
		vid := c.Param("venue_id")
		fmt.Println(vid)
		uid := c.GetString("uid")
		file, err := c.FormFile("image")

		if err != nil {
			fmt.Println("Error in uploading Image : ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var venue models.Venue

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_id, err := primitive.ObjectIDFromHex(vid) // convert params to //mongodb Hex ID
		if err != nil {
			fmt.Println("Not able to convert vid")
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := c.Bind(&venue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Venue is followinfg \n %+v \n", venue)

		validationErr := validate.Struct(venue)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		venue.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.D{{Key: "_id", Value: _id}}
		fmt.Println(_id)
		fmt.Println(uid)

		update := bson.D{{Key: "$set", Value: venue}}

		result, err := venueCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := "Venue has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No documents is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		isUpdated := c.Request.FormValue("is_removed")

		if (isUpdated == constants.TRUE) && (venue.ImageURL != "") {
			fp := constants.VenueImageDir + "/" + vid
			sp := constants.VenueImageURL + "/" + vid
			err = helper.RemoveImage(venue.ImageURL, sp, fp, c)
			if err == nil {
				imageURL := helper.GetImageURL(file, vid, fp, sp, c)
				venue.ImageURL = imageURL
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, "Venue updated successfully")
	}
}

func GetVenue() gin.HandlerFunc {
	return func(c *gin.Context) {
		vid := c.Param("venue_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var venue models.Venue
		id, err := primitive.ObjectIDFromHex(vid)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = venueCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&venue)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, venue)
	}
}

func GetVenues() gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Request.URL.RawQuery

		var venues []models.Venue

		opts := mgs.FindOption()
		// Set max limit to restrict the number of results returned
		opts.SetMaxLimit(100)
		result, err := mgs.MongoGoSearch(query, opts)
		if err != nil {
			//invalid query
			fmt.Println("Invalid query", err)
			c.JSON(http.StatusInternalServerError, venues)
			return
		}

		findOpts := options.Find()
		findOpts.SetSort(result.Sort)

		projection := bson.D{
			{Key: "name", Value: 1},
			{Key: "imageurl", Value: 1},
			{Key: "owner", Value: 1},
			{Key: "location", Value: 1},
			{Key: "ctgry", Value: 1},
			{Key: "time", Value: 1},
		}
		findOpts.SetProjection(projection)

		cur, err := venueCollection.Find(c, result.Filter, findOpts)
		if err != nil {
			fmt.Print("Error finding venues", err)
			c.JSON(http.StatusInternalServerError, venues)
			return
		}
		for cur.Next(c) {
			var venue models.Venue
			err := cur.Decode(&venue)
			if err != nil {
				c.JSON(http.StatusInternalServerError, venues)
				return
			}
			venues = append(venues, venue)
		}

		c.JSON(http.StatusOK, gin.H{"count": len(venues), "venues": venues})
	}
}
