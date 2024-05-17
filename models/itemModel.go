package models

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `form:"title" json:"title"`
	Description string             `form:"desc" json:"desc" bson:"desc"`
	ImagesURL   []string           `json:"imagesurl"`
	Seller      string             `form:"seller" json:"seller"`
	Buyer       string             `json:"buyer"`
	Location    string             `form:"location" json:"location" validate:"required"`
	Status      string             `form:"status" json:"status" validate:"required,eq=ACTIVE|eq=SOLD"`
	Condition   int                `form:"condition" json:"condition"`
	Price       int                `form:"price" json:"price" validate:"required,number"`
	Category    string             `form:"ctgry" json:"ctgry" bson:"ctgry"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	User_id     string             `json:"user_id"`
	Product_id  string             `json:"pid" bson:"pid"`
}

type ListItem struct {
	ID        primitive.ObjectID `json:"omitempty" bson:"_id,omitempty"`
	ImageURL  string             `json:"imageurl"`
	Title     string             `json:"title"`
	Condition int                `json:"condition"`
	Price     int                `json:"price"`
	Seller    string             `json:"seller"`
	Category  string             `json:"ctgry"`
	User_id   string             `json:"user_id" validate:"required"`
	Status    string             `json:"status" validate:"required,eq=ACTIVE|eq=SOLD"`
}

type ItemImages struct {
	Added   []int                   `form:"added"`
	Removed []int                   `form:"removed"`
	Updated []int                   `form:"updated"`
	Files   []*multipart.FileHeader `form:"files"`
}

type ItemInfo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `form:"title" json:"title"`
	Description string             `form:"desc" json:"desc" bson:"desc,omitempty"`
	ImagesURL   []string           `json:"imagesurl" bson:"imagesurl,omitempty"`
	Seller      string             `form:"seller" json:"seller"`
	Buyer       string             `json:"buyer" bson:"buyer,omitempty"`
	Location    string             `form:"location" json:"location" validate:"required"`
	Status      string             `form:"status" json:"status" validate:"required,eq=ACTIVE|eq=SOLD"`
	Condition   int                `form:"condition" json:"condition,omitempty"`
	Price       int                `form:"price" json:"price" validate:"required,number"`
	Category    string             `form:"ctgry" json:"ctgry" bson:"ctgry,omitempty"`
	Created_at  time.Time          `json:"created_at" bson:"created_at,omitempty"`
	Updated_at  time.Time          `json:"updated_at"`
	User_id     string             `json:"user_id" bson:"user_id,omitempty"`
	Product_id  string             `json:"pid" bson:"pid,omitempty"`
}
