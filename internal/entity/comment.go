package entity

import "time"

type Comment struct {
	ID        int       `json:"id"`
	ParentID  int       `json:"parentId"`
	Level     int       `json:"level"`
	UserID    int       `json:"userId"`
	PostID    int       `json:"postId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
