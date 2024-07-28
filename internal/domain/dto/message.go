package dto

import "time"

type Message struct {
	Id        int64     `json:"id" validate:"required"`
	Body      string    `json:"body" validate:"required"`
	CreatedAt time.Time `json:"createdAt validate:"required"`
}
