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

		file, _ := c.FormFile("image")
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

func GetBodies() gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Request.URL.RawQuery

		var bodies []models.Body

		opts := mgs.FindOption()
		// Set max limit to restrict the number of results returned
		opts.SetMaxLimit(100)
		result, err := mgs.MongoGoSearch(query, opts)
		if err != nil {
			//invalid query
			log.Print("Invalid query", err)
			c.JSON(http.StatusInternalServerError, bodies)
			return
		}

		findOpts := options.Find()
		findOpts.SetLimit(result.Limit)
		findOpts.SetSkip(result.Skip)
		findOpts.SetSort(result.Sort)

		projection := bson.D{
			{Key: "username", Value: 1},
			{Key: "name", Value: 1},
			{Key: "desc", Value: 1},
			{Key: "imageurl", Value: 1},
			{Key: "location", Value: 1},
			{Key: "type", Value: 1},
			{Key: "level", Value: 1},
			{Key: "ctgry", Value: 1},
		}
		findOpts.SetProjection(projection)

		cur, err := bodyCollection.Find(c, result.Filter, findOpts)
		if err != nil {
			log.Print("Error finding products", err)
			c.JSON(http.StatusInternalServerError, bodies)
			return
		}
		for cur.Next(c) {
			var body models.Body
			err := cur.Decode(&body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, bodies)
				return
			}
			bodies = append(bodies, body)
		}

		c.JSON(http.StatusOK, gin.H{"count": len(bodies), "bodies": bodies})
	}
}

func GetBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		bid := c.Param("body_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var body models.Body
		fmt.Println(bid)
		id, err := primitive.ObjectIDFromHex(bid)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = bodyCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&body)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, body)
	}
}

func VerfiyBody() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bid := c.Param("body_id")
		fmt.Println(bid)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var body models.Body
		id, _ := primitive.ObjectIDFromHex(bid)
		err := bodyCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&body)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.D{{Key: "_id", Value: id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}}

		result, err := bodyCollection.UpdateOne(ctx, filter, update)
		fmt.Println(body)
		fmt.Println(result.MatchedCount)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Body is approved"})

		defer cancel()
	}
}

func UpdateBody() gin.HandlerFunc {

	return func(c *gin.Context) {
		bid := c.Param("body_id")
		fmt.Println(bid)
		id, err := primitive.ObjectIDFromHex(bid)

		password := c.Request.FormValue("password")
		var foundBody models.Body
		err = bodyCollection.FindOne(c, bson.M{"_id": id}).Decode(&foundBody)

		fmt.Printf("Given password %s\nFound password %s\n", password, *foundBody.Password)
		passwordIsValid, msg := helper.VerifyPassword(password, *foundBody.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		file, err := c.FormFile("image")

		if err != nil && err.Error() != constants.NO_IMAGE_IN_FORM {
			fmt.Println("Error in uploading Image : ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var body models.Body

		_id, err := primitive.ObjectIDFromHex(bid) // convert params to //mongodb Hex ID
		if err != nil {
			fmt.Println("Not able to convert bid")
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := c.Bind(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Body is followinfg \n %+v \n", body)

		validationErr := validate.Struct(body)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		isUpdated := c.Request.FormValue("is_removed")

		if (isUpdated == constants.TRUE) && (body.ImageURL != "") {
			fp := constants.BodyLogoDir + "/" + bid
			sp := constants.BodyLogoURL + "/" + bid
			err = helper.RemoveImage(body.ImageURL, sp, fp, c)
			if err == nil {
				imageURL := helper.GetImageURL(file, bid, fp, sp, c)
				body.ImageURL = imageURL
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		body.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.D{{Key: "_id", Value: _id}}
		fmt.Println(_id)

		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "name", Value: body.Name},
			{Key: "desc", Value: body.Description},
			{Key: "address", Value: body.Address},
			{Key: "location", Value: body.Location},
			{Key: "website", Value: body.Website},
			{Key: "imageurl", Value: body.ImageURL},
		}}}

		result, err := bodyCollection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			fmt.Println(err.Error())
			msg := "Body has not been updated"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.MatchedCount == 0 {
			msg := "No documents is matched"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, "Body updated successfully")
	}
}

func PutCouncilStruct() gin.HandlerFunc {
	return func(c *gin.Context) {

		bid := c.Param("body_id")
		id, _ := primitive.ObjectIDFromHex(bid)

		var councilStruct models.Body
		if err := c.BindJSON(&councilStruct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		password := councilStruct.Password
		var foundBody models.Body
		//opts := options.FindOne().SetProjection(bson.D{{Key: "password", Value: 1}, {Key: "verified", Value: 1}})
		_ = bodyCollection.FindOne(c, bson.M{"_id": id}).Decode(&foundBody)

		passwordIsValid, msg := helper.VerifyPassword(*password, *foundBody.Password)
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if !foundBody.Verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Body is not verified"})
			return
		}

		validationErr := validate.Struct(councilStruct.Council)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.D{{Key: "_id", Value: id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "council", Value: councilStruct.Council},
			{Key: "updated_at", Value: updated_at},
		}}}

		result, err := bodyCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println(err.Error())
			msg := fmt.Sprintf("Council Struct has not been updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		fmt.Println(result.MatchedCount)

		c.JSON(http.StatusOK, "Council Struct updated successfully")

	}
}
