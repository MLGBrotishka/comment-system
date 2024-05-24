package usecase

import (
	"comment-system/internal/entity"
	"context"
)

type CommentsRepo interface {
	GetById(ctx context.Context, id int) (entity.Comment, error)
	GetReplies(ctx context.Context, id int, limit *int, offset *int) ([]entity.Comment, error)
	GetByPost(ctx context.Context, postID int, limit *int, offset *int) ([]entity.Comment, error)
	Store(ctx context.Context, new entity.NewComment) (entity.Comment, error)
}

type PostsRepo interface {
	GetById(ctx context.Context, id int) (entity.Post, error)
	GetAll(ctx context.Context, limit *int, offset *int) ([]entity.Post, error)
	UpdatePermission(ctx context.Context, postID int, enable bool) (entity.Post, error)
	Store(ctx context.Context, new entity.NewPost) (entity.Post, error)
}

type PostGetter interface {
	GetById(ctx context.Context, id int) (entity.Post, error)
}
