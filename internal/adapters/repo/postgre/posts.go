package postgre

import (
	"comment-system/internal/entity"
	"comment-system/pkg/postgres"
	"context"
	"fmt"
)

type PostsRepo struct {
	limit  int
	offset int
	pg     *postgres.Postgres
}

func NewPostsRepo(pg *postgres.Postgres, defaultLimit, defaultOffset int) *PostsRepo {
	return &PostsRepo{
		limit:  defaultLimit,
		offset: defaultOffset,
		pg:     pg,
	}
}

func (r *PostsRepo) GetById(ctx context.Context, id int) (entity.Post, error) {
	query := `SELECT p.id, p.user_id, p.comments_enabled, p.name, p.text, p.created_at FROM posts AS p WHERE p.id=$1;`
	row := r.pg.Pool.QueryRow(ctx, query, id)
	post := entity.Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.CommentsEnabled, &post.Name, &post.Text, &post.CreatedAt)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsRepo - GetById - row.Scan: %w", err)
	}
	return post, nil
}

func (r *PostsRepo) GetAll(ctx context.Context, limit *int, offset *int) ([]entity.Post, error) {
	query := `SELECT p.id, p.user_id, p.comments_enabled, p.name, p.text, p.created_at FROM posts AS p ORDER BY p.created_at DESC LIMIT $1 OFFSET $2;`
	if limit == nil {
		limit = &r.limit
	}
	if offset == nil {
		offset = &r.offset
	}
	rows, err := r.pg.Pool.Query(ctx, query, *limit, *offset)
	if err != nil {
		return []entity.Post{}, fmt.Errorf("PostsRepo - GetAll - Query: %w", err)
	}
	defer rows.Close()
	var posts []entity.Post
	for rows.Next() {
		post := entity.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.CommentsEnabled, &post.Name, &post.Text, &post.CreatedAt)
		if err != nil {
			return []entity.Post{}, fmt.Errorf("PostsRepo - GetAll - Scan: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostsRepo) UpdatePermission(ctx context.Context, postID int, enable bool) (entity.Post, error) {
	query := `UPDATE posts SET comments_enabled=$1 WHERE id=$2 RETURNING id, user_id, comments_enabled, name, text, created_at;`
	row := r.pg.Pool.QueryRow(ctx, query, enable, postID)
	post := entity.Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.CommentsEnabled, &post.Name, &post.Text, &post.CreatedAt)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsRepo - UpdatePermission - row.Scan: %w", err)
	}
	return post, nil
}

func (r *PostsRepo) Store(ctx context.Context, new entity.NewPost) (entity.Post, error) {
	query := `INSERT INTO posts (user_id, comments_enabled, name, text) VALUES ($1, $2, $3, $4) RETURNING id, user_id, comments_enabled, name, text, created_at;`
	row := r.pg.Pool.QueryRow(ctx, query, new.UserID, new.CommentsEnabled, new.Name, new.Text)
	post := entity.Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.CommentsEnabled, &post.Name, &post.Text, &post.CreatedAt)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostsRepo - Store - Scan: %w", err)
	}
	return post, nil
}
