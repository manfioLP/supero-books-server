package models
import "go.mongodb.org/mongo-driver/bson/primitive"
type ToDoList struct {

	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string             `json:"title" validate:"required,min=2,max=100"`
	ISBN   string             `json:"ISBN" validate:"required,min=13,max=13"`
	Author   string             `json:"author,omitempty" validate:"required,min=2,max=100"`
	Publisher   string             `json:"publisher" validate:"required,min=2,max=100"`
	Language string               `json:"language" validate:"required,min=2,max=100"`
	Weight int32               `json:"weight,omitempty" validate:"gte=0,lte=10000"`
	Height int32               `json:"height,omitempty" validate:"gte=0,lte=10000"`
	Length int32               `json:"length,omitempty" validate:"gte=0,lte=10000"`
	Width int32               `json:"width,omitempty" validate:"gte=0,lte=10000"`

}

