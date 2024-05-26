package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BodyLevels map[string]int = map[string]int{"DOSA": 0, "Institute": 1, "Individuals": 2}

// TODO: Add strict validation
type Body struct {
	// To control access of body admin to the features and the platform
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Password string             `form:"password" json:"password" bson:"password"`
	Verified bool               `form:"verified" json:"verified" bson:"verified"`

	//Basic details about the body
	Name        string      `form:"name" json:"name" validate:"required"`
	Description string      `form:"desc" json:"desc,omitempty" bson:"desc" validate:"required"`
	ImageURL    string      `form:"imageurl" json:"imageurl,omitempty" bson:"imageurl"`
	Address     string      `form:"address" json:"address" bson:"address"`
	Location    Coordinates `form:"location" json:"location" validate:"required"`
	Website     string      `form:"website" json:"website,omitempty" bson:"website"`

	// Definig the level, category and the type of body
	Type     string `form:"type" json:"type" bson:"type"`
	Level    int8   `form:"level" json:"level" bson:"level"`
	Category string `form:"ctgry" json:"ctgry" bson:"ctgry" validate:"required"`

	// Structure of the body
	Parent  string        `form:"pid" json:"pid" bson:"pid"`
	Council CouncilStruct `form:"council" json:"council" bson:"council"`

	// Database record
	Created_at time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_at time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CouncilStruct struct {
	Fourthie []string
	Thirdie  []string
	Sophie   []string
	Freshie  []string
}

type Member struct {
	Name string
	Uid  string
}

type Council struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Body     string             `form:"bid" json:"bid" bson:"bid"`
	Fourthie []Member
	Thirdie  []Member
	Sophie   []Member
	Freshie  []Member
}
