package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/radifan9/social-media-backend/internal/models"
	"github.com/radifan9/social-media-backend/internal/repositories"
	"github.com/radifan9/social-media-backend/internal/utils"
	"github.com/radifan9/social-media-backend/pkg"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	ur *repositories.UserRepository
	ac *utils.AuthCacheManager
}

func NewUserHandler(ur *repositories.UserRepository, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		ur: ur,
		ac: utils.NewAuthCacheManager(rdb),
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
