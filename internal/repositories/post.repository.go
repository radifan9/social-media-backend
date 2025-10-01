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

	// Step 2 : Insert images (if any)
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

func (p *PostRepository) GetFollowingFeed(ctx context.Context, userID string) ([]models.FeedPost, error) {
	query := `
		SELECT 
			p.user_id,
			p.text_content,
			p.created_at,
			up.name as author_name,
			COUNT(DISTINCT pl.id) as like_count,
			ARRAY_AGG(DISTINCT pi.image_url) FILTER (WHERE pi.image_url IS NOT NULL) as images,
			JSON_AGG(
				JSONB_BUILD_OBJECT(
					'name', COALESCE(cup.name, cu.email),
					'comment_text', pc.comment,
					'created_at', pc.created_at
				) ORDER BY pc.created_at DESC
			) FILTER (WHERE pc.id IS NOT NULL) as comments
		FROM posts p
		INNER JOIN user_followers uf ON p.user_id = uf.user_id
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON p.user_id = up.user_id
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		LEFT JOIN post_comments pc ON p.id = pc.post_id
		LEFT JOIN users cu ON pc.user_id = cu.id
		LEFT JOIN user_profiles cup ON pc.user_id = cup.user_id
		LEFT JOIN post_images pi ON p.id = pi.post_id
		WHERE uf.follower_id = $1
		GROUP BY p.id, p.user_id, p.text_content, p.created_at, u.email, up.name, up.avatar
		ORDER BY p.created_at DESC
		LIMIT 10
	`

	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return []models.FeedPost{}, err
	}
	defer rows.Close()

	var posts []models.FeedPost

	for rows.Next() {
		var post models.FeedPost
		var comments []models.FeedComment

		if err := rows.Scan(
			&post.UserID,
			&post.TextContent,
			&post.CreatedAt,
			&post.AuthorName,
			&post.LikeCount,
			&post.Images,
			&comments,
		); err != nil {
			return []models.FeedPost{}, err
		}

		if comments != nil {
			post.Comments = []models.FeedComment(comments)
		} else {
			post.Comments = []models.FeedComment{}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return []models.FeedPost{}, err
	}

	return posts, nil
}
