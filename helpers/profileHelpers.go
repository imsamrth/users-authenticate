package helper

import (
	"context"
	"errors"
	"fmt"
	models "instix_auth/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ValidateUsername(username string, user_id string, profileCollection *mongo.Collection, ctx context.Context) (statusCode int, err error) {

	var result models.Profile
	filter := bson.D{{Key: "username", Value: username}}
	err = profileCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, err
	}

	if result.ID.Hex() != user_id {
		fmt.Println(result.ID.Hex(), " Value of result that matched username")
		fmt.Println(user_id, " Value of provided user id")
		return http.StatusBadRequest, errors.New("this username is already taken")
	}
	return http.StatusOK, nil
}
