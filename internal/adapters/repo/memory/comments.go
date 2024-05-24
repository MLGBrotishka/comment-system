package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"comment-system/internal/entity"
)

type MemoryCommentsRepo struct {
	limit    int
	offset   int
	comments map[int]entity.Comment
	replies  map[int][]int
	post     map[int][]int
	mutex    sync.RWMutex
}

func NewMemoryCommentsRepo(defaultLimit, defaultOffset int) *MemoryCommentsRepo {
	return &MemoryCommentsRepo{
		limit:    defaultLimit,
		offset:   defaultOffset,
		comments: make(map[int]entity.Comment),
		replies:  make(map[int][]int),
		post:     make(map[int][]int),
		mutex:    sync.RWMutex{},
	}
}

func (r *MemoryCommentsRepo) GetById(ctx context.Context, id int) (entity.Comment, error) {
	r.mutex.RLock()
	comment, ok := r.comments[id]
	r.mutex.RUnlock()
	if !ok {
		return entity.Comment{}, errors.New("comment not found")
	}
	return comment, nil
}

func (r *MemoryCommentsRepo) GetReplies(ctx context.Context, id int, limit, offset *int) ([]entity.Comment, error) {
	if limit == nil {
		limit = &r.limit
	}
	if offset == nil {
		offset = &r.offset
	}
	r.mutex.RLock()
	repliesIds := r.replies[id]
	r.mutex.RUnlock()
	var replies []entity.Comment
	for _, id := range repliesIds {
		comment, err := r.GetById(ctx, id)
		if err != nil {
			continue
		}
		replies = append(replies, comment)
	}
	sort.SliceStable(replies, func(i, j int) bool {
		return replies[i].CreatedAt.Before(replies[j].CreatedAt)
	})
	if len(replies) <= *offset {
		return []entity.Comment{}, nil
	}

	start := *offset
	end := *offset + *limit
	if end > len(replies) {
		end = len(replies)
	}
	replies = replies[start:end]
	return replies, nil
}

func (r *MemoryCommentsRepo) GetByPost(ctx context.Context, id int, limit, offset *int) ([]entity.Comment, error) {
	if limit == nil {
		limit = &r.limit
	}
	if offset == nil {
		offset = &r.offset
	}
	r.mutex.RLock()
	postIds := r.post[id]
	r.mutex.RUnlock()
	var comments []entity.Comment
	for _, id := range postIds {
		comment, err := r.GetById(ctx, id)
		if err != nil {
			continue
		}
		comments = append(comments, comment)
	}
	sort.SliceStable(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})
	if len(comments) <= *offset {
		return []entity.Comment{}, nil
	}

	start := *offset
	end := *offset + *limit
	if end > len(comments) {
		end = len(comments)
	}
	comments = comments[start:end]
	return comments, nil
}

func (r *MemoryCommentsRepo) Store(ctx context.Context, new entity.NewComment) (entity.Comment, error) {
	if new.ParentID == nil {
		new.ParentID = &[]int{0}[0]
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	comment := entity.Comment{
		ID:        len(r.comments) + 1,
		ParentID:  *new.ParentID,
		Level:     1, // Assuming top-level comments always have Level 1
		UserID:    new.UserID,
		PostID:    new.PostID,
		Text:      new.Text,
		CreatedAt: time.Now(),
	}
	if comment.ParentID != 0 {
		r.replies[comment.ParentID] = append(r.replies[comment.ParentID], comment.ID)
	}
	r.post[comment.PostID] = append(r.post[comment.PostID], comment.ID)
	r.comments[comment.ID] = comment
	return comment, nil
}
