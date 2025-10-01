package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radifan9/social-media-backend/internal/handlers"
	"github.com/radifan9/social-media-backend/internal/middlewares"
	"github.com/radifan9/social-media-backend/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func RegisterPostRoutes(v1 *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client) {
	postRepo := repositories.NewPostRepository(db, rdb)
	postHandler := handlers.NewPostHandler(postRepo, rdb)
	verifyTokenWithBlacklist := middlewares.VerifyTokenWithBlacklist(rdb)

	post := v1.Group("/post")
	post.POST("/", verifyTokenWithBlacklist, postHandler.CreatePost)
	post.POST("/like", verifyTokenWithBlacklist, postHandler.LikePost)
	post.POST("/comment", verifyTokenWithBlacklist, postHandler.AddComment)

	feed := v1.Group("/feed")
	feed.GET("/", verifyTokenWithBlacklist, postHandler.GetFollowingFeed)
}
