package controllers

import (
	"context"
	"errors"
	"fmt"
	constants "instix_auth/constants"
	"instix_auth/database"
	helper "instix_auth/helpers"
	models "instix_auth/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var profileCollection *mongo.Collection = database.OpenCollection(database.Client, constants.PROFILEDATABASE)

func CreateProfile(uid string, roll_no string) (statusCode int, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var profile models.Profile
	profile.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	profile.User_id = uid
	profile.Roll_no = roll_no

	validationErr := validate.Struct(profile)
	if validationErr != nil {
		return http.StatusBadRequest, validationErr
	}

	count, err := profileCollection.CountDocuments(ctx, bson.M{"_id": uid})

	fmt.Print("Checking is profile already created")
	defer cancel()
	if err != nil {
		log.Panic(err)
		return http.StatusInternalServerError, err
	}

	if count > 0 {
		return http.StatusInternalServerError, errors.New("this profile is already created")
	}

	filter := bson.D{{Key: "_id", Value: uid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "roll_no", Value: roll_no}}}}
	opts := options.Update().SetUpsert(true)
	resultInsertionNumber, insertErr := profileCollection.UpdateOne(ctx, filter, update, opts)
	if insertErr != nil {
		fmt.Println("Profile item was not created")
		return http.StatusInternalServerError, insertErr
	}

	fmt.Printf("Documents matched: %v\n", resultInsertionNumber)

	defer cancel()
	return http.StatusOK, nil
}

// TODO : Turn off upsert methods
func UpdateUsername() gin.HandlerFunc {

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

func UpdateProfile() gin.HandlerFunc {

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

// TODO : Add limitations to type and size of image
func UploadAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetString("uid")

		println(user_id, " Update Image  is called")
		if user_id == "" {
			log.Println("No user id found ")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no username found"})
			return
		}

		file, err := c.FormFile("image")

		if err != nil {
			log.Println("Error in uploading Image : ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		uniqueId := uuid.New()

		filename := strings.Replace(uniqueId.String(), "-", "", -1)

		fileExt := strings.Split(file.Filename, ".")[1]

		image := fmt.Sprintf("%s.%s", filename, fileExt)

		err = c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", constants.ProfileImageDir, image))

		if err != nil {
			log.Println("Error in saving Image :", err)
			c.JSON(http.StatusBadRequest, err.Error())
		}

		imageUrl := fmt.Sprintf("%s%s/%s", domName, constants.ProfileImageURL, image)

		// Put image url to profile model
		filter := bson.D{{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "profileurl", Value: imageUrl}}}}
		opts := options.Update().SetUpsert(true)

		result, err := profileCollection.UpdateOne(context.TODO(), filter, update, opts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("Documents matched: %v\n", result.MatchedCount)
		fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

		c.JSON(http.StatusCreated, gin.H{"image_url": imageUrl})
	}
}

func GetProfiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := profileCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allprofiles []bson.M
		if err = result.All(ctx, &allprofiles); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allprofiles[0])
	}
}

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		pid := c.Param("profile_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		id, err := primitive.ObjectIDFromHex(pid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := profileCollection.FindOne(ctx, bson.M{"_id": id})
		err = result.Err()

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusBadRequest, gin.H{"error": constants.NO_PROFILE})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while getting profiles"})
			return
		} else {
			var profile models.Profile
			err = result.Decode(&profile)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing profile"})
				return
			}

			c.JSON(http.StatusOK, profile)
		}
	}
}
