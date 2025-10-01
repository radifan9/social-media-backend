package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radifan9/social-media-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type PostRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewPostRepository(db *pgxpool.Pool, rdb *redis.Client) *PostRepository {
	return &PostRepository{
		db:  db,
		rdb: rdb,
	}
}

func (p *PostRepository) CreatePost(ctx context.Context, userID string, body models.CreatePost, imagePaths []string) (models.Post, error) {
	// Begin transaction
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return models.Post{}, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("failed to rollback transaction:", rollbackErr)
			}
		}
	}()

	// Step  1 : Insert into posts table
	var post models.Post
	postQuery := `
		Insert into posts (user_id, text_content)
		values ($1, $2)
		returning id, user_id, text_content, created_at
	`

	if err = tx.QueryRow(ctx, postQuery, userID, body.TextContent).Scan(
		&post.ID, &post.UserID, &post.TextContent, &post.CreatedAt,
	); err != nil {
		return models.Post{}, err
	}

	// Insert images
	for _, path := range imagePaths {
		imgQuery := `
			Insert into post_images (post_id, image_url)
			values ($1, $2)
			returning id, post_id, image_url, created_at 	
		`

		var img models.PostImage
		if err = tx.QueryRow(ctx, imgQuery, post.ID, path).Scan(
			&img.ID, &img.PostID, &img.ImageURL, &img.CreatedAt,
		); err != nil {
			return models.Post{}, err
		}
		post.Images = append(post.Images, img.ImageURL)
	}

	// Commit
	if err = tx.Commit(ctx); err != nil {
		return models.Post{}, err
	}

	return post, nil
}
