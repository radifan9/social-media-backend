package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radifan9/social-media-backend/internal/handlers"
	"github.com/radifan9/social-media-backend/internal/middlewares"
	"github.com/radifan9/social-media-backend/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRoutes(v1 *gin.RouterGroup, db *pgxpool.Pool, rdb *redis.Client) {
	userRepo := repositories.NewUserRepository(db, rdb)
	userHandler := handlers.NewUserHandler(userRepo, rdb)
	verifyTokenWithBlacklist := middlewares.VerifyTokenWithBlacklist(rdb)

	auth := v1.Group("/auth")
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)
	auth.DELETE("/logout", verifyTokenWithBlacklist, userHandler.Logout)

	user := v1.Group("/user")
	user.Use(verifyTokenWithBlacklist)
	user.PATCH("/", userHandler.EditProfile)
	user.POST("/:targetID/follow", userHandler.FollowUser)
}
