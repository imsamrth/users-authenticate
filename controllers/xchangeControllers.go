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

		imagesURL, count := helper.GetImagesURL("images", item.ID.Hex(), c)

		if count != (constants.MaxItemImages + 1) {
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
func UpdateItemInfo() gin.HandlerFunc {

	return func(c *gin.Context) {
		pid := c.Param("product_id")
		uid := c.GetString("uid")

		var item models.ItemInfo

		_id, err := primitive.ObjectIDFromHex(pid) // convert params to //mongodb Hex ID
		if err != nil {
			fmt.Println(err.Error())
		}

		if err := c.BindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(item)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		item.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		//TODO: Handle changes in image ( remove or add)
		filter := bson.D{{Key: "_id", Value: _id}, {Key: "user_id", Value: uid}}
		fmt.Println(_id)
		fmt.Println(uid)

		update := bson.D{{Key: "$set", Value: item}}

		result, err := itemCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := "Item has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No documents is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, "Item updated successfully")
	}
}

func UpdateItemImages() gin.HandlerFunc {

	return func(c *gin.Context) {
		pid := c.Param("product_id")
		uid := c.GetString("uid")

		var ItemImages models.ItemImages

		_id, err := primitive.ObjectIDFromHex(pid) // convert params to //mongodb Hex ID
		if err != nil {
			fmt.Println(err.Error())
		}

		if err := c.Bind(&ItemImages); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(ItemImages)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//fmt.Printf("Item Images is following \n %+v \n", ItemImages)
		//ItemImages.Removed = append(ItemImages.Removed, "DEMO IMAGE HOST")
		fmt.Println(ItemImages.Removed)
		fmt.Println(len(ItemImages.Removed))
		imagesURL := helper.UpdateImagesURL(ItemImages.Removed, ItemImages.Files, pid, c)

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var updated_at time.Time
		updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.D{{Key: "_id", Value: _id}, {Key: "user_id", Value: uid}}

		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "imagesURL", Value: imagesURL},
			{Key: "updated_at", Value: updated_at},
		}}}

		result, err := itemCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := "Item Images has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No documents is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, "Item Images updated successfully")
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

func DeleteItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("product_id")
		uid := c.GetString("uid")

		_id, err := primitive.ObjectIDFromHex(pid) // convert params to //mongodb Hex ID
		if err != nil {
			fmt.Println(err.Error())
		}

		opts := options.Delete().SetCollation(&options.Collation{})

		res, err := itemCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}, {Key: "user_id", Value: uid}}, opts)
		if err != nil || (res.DeletedCount == 0) {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "No documents are matched"})
			return
		}

		fmt.Println(res.DeletedCount)
		c.JSON(http.StatusAccepted, gin.H{"msg": "Deleted successfully"})
	}
}
