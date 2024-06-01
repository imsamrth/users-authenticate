package helper

import (
	"context"
	"errors"
	"fmt"
	models "instix_auth/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IsAlreadyCouncilCreated(session string, bid string, collection *mongo.Collection) (err error) {
	err = nil
	count, err := collection.CountDocuments(context.TODO(), bson.D{{Key: "bid", Value: bid}, {Key: "session", Value: session}})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("council already created")
	}
	return
}

func VerifyBodyId(bid string, providedPassword string, collection *mongo.Collection) (err error) {

	err = nil
	id, _ := primitive.ObjectIDFromHex(bid)
	var foundBody models.Body
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: 1}})).Decode(&foundBody)
	if err != nil {
		err = errors.New(err.Error())
		return
	}

	passwordIsValid, msg := VerifyPassword(providedPassword, *foundBody.Password)
	if !passwordIsValid {
		err = errors.New(msg)
	}
	return err
}

func IsMemberOfBody(mid string, bid string, collection *mongo.Collection) (statusCode int, err error) {
	err = nil
	_id, _ := primitive.ObjectIDFromHex(mid)
	fmt.Printf("BID is %s\nMID is %s\n", bid, _id)
	result := collection.FindOne(context.TODO(), bson.D{{Key: "bid", Value: bid}, {Key: "_id", Value: _id}})
	err = result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}
