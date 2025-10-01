package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radifan9/social-media-backend/internal/models"
	"github.com/radifan9/social-media-backend/internal/repositories"
	"github.com/radifan9/social-media-backend/internal/utils"
	"github.com/radifan9/social-media-backend/pkg"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	ur *repositories.UserRepository
	ac *repositories.AuthCacheManager
}

func NewUserHandler(ur *repositories.UserRepository, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		ur: ur,
		ac: repositories.NewAuthCacheManager(rdb),
	}
}

// @Summary Register a new user
// @Tags    Auth
// @Accept  json
// @Produce json
// @Param   body body models.RegisterUser true "User registration"
// @Success 201 {object} models.User
// @Router  /api/v1/auth/register [post]
func (u *UserHandler) Register(ctx *gin.Context) {
	var user models.RegisterUser
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	log.Println("email : ", user.Email)
	log.Println("password : ", user.Password)

	// Hash password
	// "password": "ceganssangar123(DF&&"
	// format : email + sangar123(DF&&
	hashCfg := pkg.NewHashConfig()
	hashCfg.UseRecommended()
	hashedPassword, err := hashCfg.GenHash(user.Password)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to hash password", err.Error())
		return
	}

	newUser, err := u.ur.CreateUser(ctx, user.Email, hashedPassword)
	if err != nil {
		log.Println("error : ", err)
		utils.HandleError(ctx, http.StatusConflict, "failed to register", err.Error())
		return
	}

	utils.HandleResponse(ctx, http.StatusOK, models.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Data: gin.H{
			"id":    newUser.Id,
			"email": newUser.Email,
		},
	})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "bad request", err)
		// utils.HandleError(ctx, http.StatusBadRequest, "bad request", err.Error())
		return
	}

	// GetID from Database
	infoUser, err := u.ur.GetIDFromEmail(ctx, user.Email)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "bad request", err.Error())
		return
	}

	// Get password & role from where ID is match
	userCred, err := u.ur.GetPasswordFromID(ctx, infoUser.Id)
	if err != nil {
		log.Println("error getting password & role")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Bandingkan password
	hashCfg := pkg.NewHashConfig()
	isMatched, err := hashCfg.CompareHashAndPassword(user.Password, userCred.Password)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		re := regexp.MustCompile("hash|crypto|argon2id|format")
		if re.Match([]byte(err.Error())) {
			log.Println("Error during Hashing")
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	if !isMatched {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Nama atau Password salah",
		})
		return
	}

	// Jika match, maka buatkan jwt dan kirim via response
	claims := pkg.NewJWTClaims(infoUser.Id)
	jwtToken, err := claims.GenToken()
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	utils.Success(ctx, http.StatusOK, models.SuccessLoginResponse{
		Token: jwtToken,
	})

}

func (u *UserHandler) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		utils.Error(ctx, http.StatusBadRequest, "bad request", fmt.Errorf("authorization header is required"))
		return
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		utils.HandleError(ctx, http.StatusBadRequest, "invalid authorization format", "authorization header must be in format 'Bearer <token>'")
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		utils.HandleError(ctx, http.StatusUnauthorized, "unauthorized", "token claims not found")
		return
	}

	userClaims, ok := claims.(pkg.Claims)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "internal server error", "cannot cast claims")
		return
	}

	// Calculate remaining TTL for the token
	expirationTime := time.Unix(userClaims.ExpiresAt.Unix(), 0)
	remainingTTL := time.Until(expirationTime)

	log.Println("expirationTime : ", expirationTime)
	log.Println("remainingTTL : ", remainingTTL)

	// Only blacklist if token hasn't expired yet
	if remainingTTL > 0 {
		if err := u.ac.BlacklistToken(ctx.Request.Context(), tokenString, remainingTTL); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "internal server error", "failed to logout")
			return
		}
	}

	utils.HandleResponse(ctx, http.StatusOK, models.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Data: map[string]interface{}{
			"message": "Logout successful",
		},
	})
}

func (u *UserHandler) EditProfile(ctx *gin.Context) {
	// Get image from form-data
	var body models.EditUserProfile
	if err := ctx.ShouldBind(&body); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	// Get the userID from token
	claims, _ := ctx.Get("claims")
	user, ok := claims.(pkg.Claims)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "internal server error", "cannot cast into pkg.claims")
		return
	}

	// Dari postman harus ambil gambar baru
	file := body.Avatar
	if file != nil {
		ext := filepath.Ext(file.Filename)
		re := regexp.MustCompile(`(?i)\.(png|jpg|jpeg|webp)$`)
		if !re.MatchString(ext) {
			utils.HandleError(ctx, http.StatusBadRequest, "invalid file type", "only png, jpg, jpeg, webp allowed")
			return
		}

		filename := fmt.Sprintf("%d_images_%s%s", time.Now().UnixNano(), user.UserId, ext)
		location := filepath.Join("public/avatars", filename)

		if err := ctx.SaveUploadedFile(file, location); err != nil {
			utils.HandleError(ctx, http.StatusBadRequest, err.Error(), "failed to upload")
			return
		}

		// Update profile with new image
		editedProfile, err := u.ur.EditProfile(ctx.Request.Context(), user.UserId, body, filename)
		if err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, err.Error(), "cannot edit user profile")
			return
		}

		utils.HandleResponse(ctx, http.StatusOK, models.SuccessResponse{Success: true, Status: http.StatusOK, Data: editedProfile})
		return
	}

	// If no image uploaded, just update profile without image
	editedProfile, err := u.ur.EditProfile(ctx.Request.Context(), user.UserId, body, "")
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "internal server error", err)
		return
	}

	utils.Success(ctx, http.StatusOK, editedProfile)
}

func (u *UserHandler) FollowUser(ctx *gin.Context) {
	// Get the userID from token
	claims, _ := ctx.Get("claims")
	user, ok := claims.(pkg.Claims)
	if !ok {
		utils.HandleError(ctx, http.StatusInternalServerError, "internal server error", "cannot cast into pkg.claims")
		return
	}

	// Get follow target
	targetID := ctx.Param("targetID")

	if err := u.ur.FollowUser(ctx, user.UserId, targetID); err != nil {
		switch {
		case errors.Is(err, repositories.ErrAlreadyFollowed):
			utils.Error(ctx, http.StatusConflict, "you already follow this user.", err)
		default:
			utils.Error(ctx, http.StatusInternalServerError, "internal server error", err)
		}
		return
	}

	utils.Success(ctx, http.StatusOK, nil)
}
