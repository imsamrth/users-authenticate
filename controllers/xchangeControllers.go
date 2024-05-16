package controllers

import (
	"context"
	"fmt"
	constants "instix_auth/constants"
	"instix_auth/database"
	helper "instix_auth/helpers"
	models "instix_auth/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var itemCollection *mongo.Collection = database.OpenCollection(database.Client, constants.ITEMDATABASE)

func CreateItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var item models.Item

		defer cancel()
		if err := c.Bind(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Form binded successfully")

		validationErr := validate.Struct(item)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		fmt.Println("Validation checked successfully")

		item.User_id = c.GetString("uid")
		item.Seller = c.GetString("first_name")
		item.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		item.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		item.ID = primitive.NewObjectID()
		item.Product_id = item.ID.Hex()

		imagesURL, count := helper.GetImageURL("images", item.ID.Hex(), c)

		if count != 10 {
			msg := "Error uploading all the images"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		item.ImagesURL = imagesURL

		resultInsertionNumber, insertErr := itemCollection.InsertOne(ctx, item)
		if insertErr != nil {
			msg := "Item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// TODO : Turn off upsert methods
func UpdateItem() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.GetString("uid")

		println(user_id, " Update Profile is called")
		if user_id == "" {
			log.Println("No user id found ")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no username found"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var profile models.Profile
		defer cancel()
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			fmt.Println(profile.Username)
		}
		validationErr := validate.Struct(profile)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		statusCode, err := helper.ValidateUsername(*&profile.Username, user_id, profileCollection, ctx)

		if err != nil {
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}
		defer cancel()

		profile.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.D{{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "username", Value: profile.Username},
			{Key: "updated_at", Value: profile.Updated_at},
		}}}

		_, err = profileCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := fmt.Sprintf("Profile has not been updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Username updated successfully")
	}
}

func UpdateItems() gin.HandlerFunc {

	return func(c *gin.Context) {
		user_id := c.GetString("uid")
		isPrimary := c.Param("isPrimary")
		println(isPrimary, "Value of isPrimary")
		println(user_id, " Update Profile is called")
		if user_id == "" {
			log.Println("No user id found ")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no username found"})
			return
		}

		var profile models.Profile
		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(profile)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		profile.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.D{{Key: "_id", Value: user_id}}
		var update bson.D
		if isPrimary == "true" {
			update = bson.D{{Key: "$set", Value: bson.D{
				{Key: "updated_at", Value: profile.Updated_at},
				{Key: "hostel", Value: profile.Hostel},
				{Key: "department", Value: profile.Department},
				{Key: "age", Value: profile.Age},
				{Key: "clubs", Value: profile.Clubs},
				{Key: "teams", Value: profile.Teams},
				{Key: "year", Value: profile.Batch},
			}}}
		} else {
			update = bson.D{{Key: "$set", Value: bson.D{
				{Key: "updated_at", Value: profile.Updated_at},
				{Key: "interest", Value: profile.Interest},
				{Key: "por", Value: profile.POR},
				{Key: "social", Value: profile.Social},
				{Key: "relationship", Value: profile.Relationship},
			}}}
		}
		opts := options.Update().SetUpsert(true)

		result, err := profileCollection.UpdateOne(context.TODO(), filter, update, opts)

		if err != nil {
			fmt.Println(err.Error())
			msg := fmt.Sprintf("Profile has not been updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result.ModifiedCount)
	}
}

func GetItems() gin.HandlerFunc {
	return func(c *gin.Context) {

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
				{"product_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := itemCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing products"})
		}
		var allProducts []bson.M
		if err = result.All(ctx, &allProducts); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allProducts[0])
	}
}

func GetItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("product_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var item models.Item
		fmt.Println(pid)
		id, err := primitive.ObjectIDFromHex(pid)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = itemCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func DeleteItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}
