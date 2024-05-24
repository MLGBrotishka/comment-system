package entity

type NewPost struct {
	UserID          int    `json:"userId"`
	CommentsEnabled *bool  `json:"commentsEnabled,omitempty"`
	Name            string `json:"name"`
	Text            string `json:"text"`
}
