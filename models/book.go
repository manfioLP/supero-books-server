package models
type Book struct {

	ID     string `json:"_id,omitempty" bson:"_id,omitempty"`
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

type Answer struct {
	Ok bool `json:"ok"`
	Errors []string `json:"errors,omitempty"`
	Data *Book `json:"data,omitempty"`
}

type GetAnswer struct {
	Ok bool `json:"ok"`
	Errors []string `json:"errors,omitempty"`
	Data []Book `json:"data,omitempty"`
	Page int64 `json:"page"`
	PerPage int64 `json:"perPage"`
	Total int64 `json:"total"`
}

