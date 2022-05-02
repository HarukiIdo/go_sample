package service

import "database/sql"

type CommentService struct {
	db *sql.DB
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{
		db: db,
	}
}
