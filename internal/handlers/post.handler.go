package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radifan9/social-media-backend/internal/models"
	"github.com/radifan9/social-media-backend/internal/repositories"
	"github.com/radifan9/social-media-backend/internal/utils"
	"github.com/radifan9/social-media-backend/pkg"
	"github.com/redis/go-redis/v9"
)

type PostHandler struct {
	pr *repositories.PostRepository
	ac *repositories.AuthCacheManager
}

func NewPostHandler(pr *repositories.PostRepository, rdb *redis.Client) *PostHandler {
	return &PostHandler{
		pr: pr,
		ac: repositories.NewAuthCacheManager(rdb),
	}
}

func (p *PostHandler) CreatePost(ctx *gin.Context) {
	// Get image from form-data
	var body models.CreatePost
	if err := ctx.ShouldBind(&body); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "internal server error", err)
		return
	}

	// Get the userID from token
	claims, _ := ctx.Get("claims")
	user, ok := claims.(pkg.Claims)
	if !ok {
		utils.Error(ctx, http.StatusInternalServerError, "internal server error", errors.New("cannot cast into pkg.claims"))
		return
	}

	// Get Images if Exists
	var imagePaths []string
	if len(body.Images) > 0 {
		for _, file := range body.Images {
			if file == nil {
				continue
			}

			// Validate extension
			ext := filepath.Ext(file.Filename)
			re := regexp.MustCompile(`(?i)\.(png|jpg|jpeg|webp)$`)
			if !re.MatchString(ext) {
				utils.HandleError(ctx, http.StatusBadRequest, "invalid file type", "only png, jpg, jpeg, webp allowed")
				return
			}

			// Generate unique filename
			filename := fmt.Sprintf("%d_images_%s%s", time.Now().UnixNano(), user.UserId, ext)
			location := filepath.Join("public/post_images", filename)

			// Save file
			if err := ctx.SaveUploadedFile(file, location); err != nil {
				utils.HandleError(ctx, http.StatusBadRequest, err.Error(), "failed to upload")
				return
			}

			imagePaths = append(imagePaths, filename)
		}
	}

	post, err := p.pr.CreatePost(ctx, user.UserId, body, imagePaths)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "failed to create a post", err)
		return
	}

	utils.Success(ctx, http.StatusOK, post)
}

func (p *PostHandler) GetFollowingFeed(ctx *gin.Context) {
	// Get the userID from token
	claims, _ := ctx.Get("claims")
	user, ok := claims.(pkg.Claims)
	if !ok {
		utils.Error(ctx, http.StatusInternalServerError, "internal server error", errors.New("cannot cast into pkg.claims"))
		return
	}

	posts, err := p.pr.GetFollowingFeed(ctx, user.UserId)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "internal server error", err)
		return
	}

	utils.Success(ctx, http.StatusOK, posts)
}

func (p *PostHandler) LikePost(ctx *gin.Context) {
	// Get the userID from token
	claims, _ := ctx.Get("claims")
	user, ok := claims.(pkg.Claims)
	if !ok {
		utils.Error(ctx, http.StatusInternalServerError, "internal server error", errors.New("cannot cast into pkg.claims"))
		return
	}

	// Bind request body
	var body models.LikeRequest
	if err := ctx.ShouldBind(&body); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	// Like the post
	like, err := p.pr.LikePost(ctx, user.UserId, body.PostID)
	if err != nil {
		if err.Error() == "post not found" {
			utils.HandleError(ctx, http.StatusNotFound, "post not found", err.Error())
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "failed to like post", err)
		return
	}

	utils.Success(ctx, http.StatusOK, like)
}
