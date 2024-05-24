package usecase

import (
	"comment-system/internal/entity"
	"context"
	"fmt"
	"math/rand"
)

type CommentsUseCase struct {
	repo        CommentsRepo
	commChannel map[int]map[string]chan entity.Comment
}

func NewComment(repo CommentsRepo) *CommentsUseCase {
	uc := &CommentsUseCase{
		repo:        repo,
		commChannel: make(map[int]map[string]chan entity.Comment),
	}
	return uc
}

func (u *CommentsUseCase) Create(ctx context.Context, input entity.NewComment) (entity.Comment, error) {
	comment, err := u.repo.Store(ctx, input)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentsUseCase - Create - u.repo.Store: %w", err)
	}
	mapChan, ok := u.commChannel[comment.ID]
	if !ok {
		return comment, nil
	}
	for _, ch := range mapChan {
		ch <- comment
	}
	return comment, nil
}

func (u *CommentsUseCase) GetToPost(ctx context.Context, postID int, limit *int, offset *int) ([]entity.Comment, error) {
	comments, err := u.repo.GetByPost(ctx, postID, limit, offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("CommentsUseCase - GetToPost - u.repo.GetByPost: %w", err)
	}
	return comments, nil
}

func (u *CommentsUseCase) GetReplies(ctx context.Context, commentID int, limit *int, offset *int) ([]entity.Comment, error) {
	comments, err := u.repo.GetReplies(ctx, commentID, limit, offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("CommentsUseCase - GetReplies - u.repo.GetReplies: %w", err)
	}
	return comments, nil
}

func (u *CommentsUseCase) Listen(ctx context.Context, postID int) (<-chan entity.Comment, error) {
	ch := make(chan entity.Comment, 1)
	mapChan, ok := u.commChannel[postID]
	if !ok {
		u.commChannel[postID] = make(map[string]chan entity.Comment)
		mapChan = u.commChannel[postID]
	}
	id := randStringBytes(10)
	go func() {
		<-ctx.Done()
		defer delete(u.commChannel[postID], id)
		defer close(ch)
	}()
	mapChan[id] = ch
	return ch, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
