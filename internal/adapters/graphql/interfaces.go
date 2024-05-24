package graphql

import (
	"comment-system/internal/entity"
	"context"
)

type CommentsUseCase interface {
	Create(ctx context.Context, input entity.NewComment) (entity.Comment, error)
	GetToPost(ctx context.Context, postID int, limit *int, offset *int) ([]entity.Comment, error)
	GetReplies(ctx context.Context, commentID int, limit *int, offset *int) ([]entity.Comment, error)
	Listen(ctx context.Context, postID int) (<-chan entity.Comment, error)
}

type PostsUseCase interface {
	Create(ctx context.Context, newPost entity.NewPost) (entity.Post, error)
	UpdatePermission(ctx context.Context, postId, userId int, enable bool) (entity.Post, error)
	GetAll(ctx context.Context, limit *int, offset *int) ([]entity.Post, error)
	GetById(ctx context.Context, id int) (entity.Post, error)
}
