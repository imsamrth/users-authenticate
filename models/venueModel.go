package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Venue struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `form:"name" json:"name" validate:"required"`
	Description string             `form:"desc" json:"desc" bson:"desc" validate:"required"`
	ImageURL    string             `form:"imageurl" json:"imageurl,omitempty" bson:"imageurl"`
	Owner       string             `form:"owner" json:"owner"`
	OwnerId     string             `form:"owner_id" json:"owner_id"`
	Location    Coordinates        `form:"location" json:"location" validate:"required"`
	Category    string             `form:"ctgry" json:"ctgry" bson:"ctgry" validate:"required,eq=S|eq=E|eq=B|eq=T|eq=C|eq=D|eq=L"`
	Time        Timing             `form:"time" json:"time,omitempty" bson:"time"`
	Website     string             `json:"website,omitempty" bson:"website"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Creator     string             `json:"creator" bson:"creator"`
}

type Coordinates struct {
	Lat float32 `json:"lat" bson:"lat"`
	Lon float32 `json:"lon" bson:"lon"`
}

type Timing struct {
	Open  int16 `json:"open" bson:"open"`
	Close int16 `json:"close" bson:"close"`
}
