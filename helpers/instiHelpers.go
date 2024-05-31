package helper

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
