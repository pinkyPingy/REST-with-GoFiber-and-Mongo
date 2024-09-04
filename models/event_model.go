package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Description string             `json:"description"`
	Date        string             `json:"date,omitempty" validate:"required"`
	Time        string             `json:"time,omitempty" validate:"required"`
	Location    string             `json:"location,omitempty" validate:"required"`
	Amount      string             `json:"amount,omitempty" validate:"required"`
}
