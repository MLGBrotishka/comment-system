package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"comment-system/internal/entity"
)

type MemoryPostsRepo struct {
	limit  int
	offset int
	posts  map[int]entity.Post
	mutex  sync.RWMutex
}

func NewMemoryPostsRepo(defaultLimit, defaultOffset int) *MemoryPostsRepo {
	return &MemoryPostsRepo{
		limit:  defaultLimit,
		offset: defaultOffset,
		posts:  make(map[int]entity.Post),
		mutex:  sync.RWMutex{},
	}
}

func (r *MemoryPostsRepo) GetById(ctx context.Context, id int) (entity.Post, error) {
	r.mutex.RLock()
	post, ok := r.posts[id]
	r.mutex.RUnlock()
	if !ok {
		return entity.Post{}, errors.New("post not found")
	}
	return post, nil
}

func (r *MemoryPostsRepo) GetAll(ctx context.Context, limit *int, offset *int) ([]entity.Post, error) {
	if limit == nil {
		limit = &r.limit
	}
	if offset == nil {
		offset = &r.offset
	}
	var posts []entity.Post
	r.mutex.RLock()
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	r.mutex.RUnlock()
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})

	if len(posts) <= *offset {
		return []entity.Post{}, nil
	}

	start := *offset
	end := *offset + *limit
	if end > len(posts) {
		end = len(posts)
	}
	posts = posts[start:end]
	return posts, nil
}

func (r *MemoryPostsRepo) UpdatePermission(ctx context.Context, postID int, enable bool) (entity.Post, error) {
	r.mutex.RLock()
	post, ok := r.posts[postID]
	r.mutex.RUnlock()
	if !ok {
		return entity.Post{}, errors.New("post not found")
	}
	post.CommentsEnabled = enable
	r.mutex.Lock()
	r.posts[postID] = post
	r.mutex.Unlock()
	return post, nil
}

func (r *MemoryPostsRepo) Store(ctx context.Context, new entity.NewPost) (entity.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	post := entity.Post{
		ID:              len(r.posts) + 1,
		UserID:          new.UserID,
		CommentsEnabled: *new.CommentsEnabled,
		Name:            new.Name,
		Text:            new.Text,
		CreatedAt:       time.Now(),
	}
	r.posts[post.ID] = post
	return post, nil
}
