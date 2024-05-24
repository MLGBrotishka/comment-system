package entity

type NewComment struct {
	ParentID *int   `json:"parentId,omitempty"`
	UserID   int    `json:"userId"`
	PostID   int    `json:"postId"`
	Text     string `json:"text"`
}
