package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Profile_URL  string             `json:"profileurl"`
	Username     string             `json:"username"`
	Roll_no      string             `json:"roll_no"`
	Updated_at   time.Time          `json:"updated_at"`
	User_id      string             `json:"user_id"`
	Hostel       string             `json:"hostel"`
	Department   string             `json:"department"`
	Batch        string             `json:"batch"`
	Age          int                `json:"age"`
	Clubs        []string           `json:"clubs"`
	Teams        []string           `json:"teams"`
	Interest     []string           `json:"interest"`
	POR          map[string]string  `json:"por"`
	Social       map[string]string  `json:"social"`
	Relationship string             `json:"relationship"`
}
