package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	ImagesURL   []string           `json:"imagesurl"`
	Title       string             `form:"title" json:"title"`
	Description string             `form:"desc" json:"desc"`
	Seller      string             `form:"seller" json:"seller"`
	Buyer       string             `json:"buyer"`
	Status      string             `form:"status" json:"status" validate:"required,eq=ACTIVE|eq=SOLD"`
	Condition   int                `form:"condition" json:"condition"`
	Price       int                `form:"price" json:"price" validate:"required,number"`
	Category    string             `form:"ctgry" json:"ctgry"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	User_id     string             `json:"user_id"`
	Product_id  string             `json:"pid" bson:"pid"`
}

type ListItem struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ImageURL  string             `json:"imageurl"`
	Title     string             `json:"title"`
	Condition int                `json:"condition"`
	Price     int                `json:"price"`
	Seller    string             `json:"seller"`
	Category  string             `json:"ctgry"`
	User_id   string             `json:"user_id" validate:"required"`
	Status    string             `json:"status" validate:"required,eq=ACTIVE|eq=SOLD"`
}
