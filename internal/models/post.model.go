package models

import (
	"mime/multipart"
	"time"
)

type CreatePost struct {
	TextContent string                  `form:"text-content"`
	Images      []*multipart.FileHeader `form:"images"`
}

type Post struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TextContent string    `json:"text_content"`
	CreatedAt   time.Time `json:"created_at"`
	Images      []string  `json:"images"`
}

type PostImage struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type FeedPost struct {
	UserID      string        `json:"user_id"`
	TextContent string        `json:"text_content"`
	CreatedAt   time.Time     `json:"created_at"`
	AuthorName  *string       `json:"author_name"`
	LikeCount   int           `json:"like_count"`
	Images      []string      `json:"images"`
	Comments    []FeedComment `json:"comments"`
}

type FeedComment struct {
	Name        string    `json:"name"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}
