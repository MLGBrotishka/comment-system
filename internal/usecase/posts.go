package usecase

import (
	"comment-system/internal/entity"
	"context"
	"fmt"
)

type PostsUseCase struct {
	repo PostsRepo
}

func NewPosts(repo PostsRepo) *PostsUseCase {
	uc := &PostsUseCase{
		repo: repo,
	}
	return uc
}

func (u *PostsUseCase) Create(ctx context.Context, newPost entity.NewPost) (entity.Post, error) {
	if newPost.CommentsEnabled == nil {
		newPost.CommentsEnabled = &[]bool{true}[0]
	}
	post, err := u.repo.Store(ctx, newPost)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsUseCase - Create - u.repo.Store: %w", err)
	}
	return post, nil
}

func (u *PostsUseCase) UpdatePermission(ctx context.Context, postId, userId int, enable bool) (entity.Post, error) {
	post, err := u.repo.GetById(ctx, postId)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsUseCase - UpdatePermission - u.repo.GetById: %w", err)
	}
	if post.UserID != userId {
		return entity.Post{}, fmt.Errorf("PostsUseCase - UpdatePermission - post.UserID != userId")
	}
	updatedPost, err := u.repo.UpdatePermission(ctx, postId, enable)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsUseCase - UpdatePermission - u.repo.UpdatePermission: %w", err)
	}
	return updatedPost, nil
}

func (u *PostsUseCase) GetAll(ctx context.Context, limit *int, offset *int) ([]entity.Post, error) {
	posts, err := u.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return []entity.Post{}, fmt.Errorf("PostsUseCase - GetAll - u.repo.GetAll: %w", err)
	}
	return posts, nil
}

func (u *PostsUseCase) GetById(ctx context.Context, id int) (entity.Post, error) {
	post, err := u.repo.GetById(ctx, id)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsUseCase - GetById - u.repo.GetById: %w", err)
	}
	return post, nil
}
