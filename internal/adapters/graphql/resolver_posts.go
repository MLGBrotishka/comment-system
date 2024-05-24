package graphql

import (
	"comment-system/internal/adapters/graphql/model"
	"comment-system/internal/entity"
	"context"
	"fmt"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input entity.NewPost) (entity.Post, error) {
	post, err := r.posts.Create(ctx, input)
	if err != nil {
		return entity.Post{}, fmt.Errorf("mutationResolver - CreatePost - r.posts.Create: %w", err)
	}
	return post, nil
}

// UpdatePostPermission is the resolver for the updatePostPermission field.
func (r *mutationResolver) UpdatePostPermission(ctx context.Context, input model.UpdatePostPermission) (entity.Post, error) {
	post, err := r.posts.UpdatePermission(ctx, input.ID, input.UserID, input.CommentsEnabled)
	if err != nil {
		return entity.Post{}, fmt.Errorf("mutationResolver - UpdatePostPermission - r.posts.UpdatePermission: %w", err)
	}
	return post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]entity.Post, error) {
	posts, err := r.posts.GetAll(ctx, limit, offset)
	if err != nil {
		return []entity.Post{}, fmt.Errorf("queryResolver - Posts - r.posts.GetAll: %w", err)
	}
	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id int) (entity.Post, error) {
	post, err := r.posts.GetById(ctx, id)
	if err != nil {
		return entity.Post{}, fmt.Errorf("queryResolver - Post - r.posts.GetById: %w", err)
	}
	return post, nil
}
