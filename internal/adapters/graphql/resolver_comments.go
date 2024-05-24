package graphql

import (
	"comment-system/internal/entity"
	"context"
	"fmt"
)

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *entity.Comment, limit *int, offset *int) ([]entity.Comment, error) {
	if obj == nil {
		return []entity.Comment{}, nil
	}
	comments, err := r.comments.GetReplies(ctx, obj.ID, limit, offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("commentResolver - Replies - r.comments.GetReplies: %w", err)
	}
	return comments, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input entity.NewComment) (entity.Comment, error) {
	comment, err := r.comments.Create(ctx, input)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("mutationResolver - CreateComment - r.comments.Create: %w", err)
	}
	return comment, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *entity.Post, limit *int, offset *int) ([]entity.Comment, error) {
	if obj == nil {
		return []entity.Comment{}, nil
	}
	comments, err := r.comments.GetToPost(ctx, obj.ID, limit, offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("postResolver - Comments - r.comments.GetToPost: %w", err)
	}
	return comments, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID int, limit *int, offset *int) ([]entity.Comment, error) {
	comments, err := r.comments.GetToPost(ctx, postID, limit, offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("postResolver - Comments - r.comments.GetToPost: %w", err)
	}
	return comments, nil
}

// RepliesToComment is the resolver for the repliesToComment field.
func (r *queryResolver) RepliesToComment(ctx context.Context, commentID int, limit *int, offset *int) ([]entity.Comment, error) {
	comments, err := r.comments.GetReplies(ctx, commentID, limit, offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("postResolver - Comments - r.comments.GetReplies: %w", err)
	}
	return comments, nil
}

// ListenComments is the resolver for the listenComments field.
func (r *subscriptionResolver) ListenComments(ctx context.Context, postID int) (<-chan entity.Comment, error) {
	ch, err := r.comments.Listen(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("subscriptionResolver - ListenComments - r.comments.Listen: %w", err)
	}
	return ch, nil
}
