package entity

import "time"

type Post struct {
	ID              int       `json:"id"`
	UserID          int       `json:"userId"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	Name            string    `json:"name"`
	Text            string    `json:"text"`
	CreatedAt       time.Time `json:"createdAt"`
}
