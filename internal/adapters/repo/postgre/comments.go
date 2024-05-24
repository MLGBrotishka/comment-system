package postgre

import (
	"comment-system/internal/entity"
	"comment-system/pkg/postgres"
	"context"
	"fmt"
)

type CommentsRepo struct {
	limit  int
	offset int
	pg     *postgres.Postgres
}

func NewComments(pg *postgres.Postgres, defaultLimit, defaultOffset int) *CommentsRepo {
	return &CommentsRepo{
		limit:  defaultLimit,
		offset: defaultOffset,
		pg:     pg,
	}
}

func (r *CommentsRepo) GetById(ctx context.Context, id int) (entity.Comment, error) {
	query := `SELECT c.id, c.parent_id, c.level, c.user_id, c.post_id, c.text, c.created_at FROM comments AS c WHERE c.id=$1;`
	row := r.pg.Pool.QueryRow(ctx, query, id)
	comment := entity.Comment{}
	err := row.Scan(&comment)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentsRepo - GetById - row.Scan: %w", err)
	}
	return comment, nil
}

func (r *CommentsRepo) GetReplies(ctx context.Context, id int, limit *int, offset *int) ([]entity.Comment, error) {
	query := `SELECT c.id, c.parent_id, c.level, c.user_id, c.post_id, c.text, c.created_at FROM comments 
	AS c WHERE c.parent_id=$1 ORDER BY c.created_at DESC LIMIT $2 OFFSET $3;`
	if limit == nil {
		limit = &r.limit
	}
	if offset == nil {
		offset = &r.offset
	}
	rows, err := r.pg.Pool.Query(ctx, query, id, *limit, *offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("CommentsRepo - GetReplies - Query: %w", err)
	}
	defer rows.Close()
	var comments []entity.Comment
	for rows.Next() {
		comment := entity.Comment{}
		err := rows.Scan(&comment)
		if err != nil {
			return []entity.Comment{}, fmt.Errorf("CommentsRepo - GetReplies - Scan: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentsRepo) GetByPost(ctx context.Context, postID int, limit *int, offset *int) ([]entity.Comment, error) {
	query := `SELECT c.id, c.parent_id, c.level, c.user_id, c.post_id, c.text, c.created_at FROM comments
	AS c WHERE c.post_id=$1 AND c.parent_id IS NULL ORDER BY c.created_at DESC LIMIT $2 OFFSET $3;`
	if limit == nil {
		limit = &r.limit
	}
	if offset == nil {
		offset = &r.offset
	}
	rows, err := r.pg.Pool.Query(ctx, query, postID, *limit, *offset)
	if err != nil {
		return []entity.Comment{}, fmt.Errorf("CommentsRepo - GetReplies - Query: %w", err)
	}
	defer rows.Close()
	var comments []entity.Comment
	for rows.Next() {
		comment := entity.Comment{}
		err := rows.Scan(&comment)
		if err != nil {
			return []entity.Comment{}, fmt.Errorf("CommentsRepo - GetReplies - Scan: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentsRepo) Store(ctx context.Context, new entity.NewComment) (entity.Comment, error) {
	query := `INSERT INTO comments (parent_id, level, user_id, post_id, text) VALUES ($1, $2, $3, $4, $5) RETURNING id, parent_id, level, user_id, post_id, text, created_at;`
	level := 1
	if new.ParentID != nil {
		parrent, err := r.GetById(ctx, *new.ParentID)
		if err != nil {
			return entity.Comment{}, fmt.Errorf("CommentsRepo - Store - GetById: %w", err)
		}
		level = parrent.Level + 1
	}
	row := r.pg.Pool.QueryRow(ctx, query, new.ParentID, level, new.UserID, new.PostID, new.Text)
	comment := entity.Comment{}
	err := row.Scan(&comment)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentsRepo - Store - Scan: %w", err)
	}
	return comment, nil
}
