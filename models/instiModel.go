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
	Password *string            `form:"password" json:"password" bson:"password"`
	Verified bool               `form:"verified" json:"verified" bson:"verified"`
	Username string             `form:"username" json:"username" bson:"username"`
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
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	POR     string
	Body    string
	BID     string `json:"bid" bson:"bid"`
	Session string `json:"session" bson:"session"`
	Level   int8
	Name    string
	Uid     string
	Tags    []string
}

type Council struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Body     string             `form:"bid" json:"bid" bson:"bid"`
	Session  string
	Fourthie []Member
	Thirdie  []Member
	Sophie   []Member
	Freshie  []Member
}

type Event struct {
	Name      string
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Desc      string
	Image     string             `json:"image" bson:"image"`
	BID       string             `json:"bid" bson:"bid"`
	Type      string             `json:"type" bson:"type"`
	Vid       string             `json:"vid" bson:"vid"`
	StartTime primitive.DateTime `json:"start_time" bson:"start_time" validate:"required"`
	EndTime   primitive.DateTime `json:"end_time" bson:"end_time" validate:"required"`
	Tags      []string
}
