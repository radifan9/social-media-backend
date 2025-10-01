package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/radifan9/social-media-backend/internal/models"
)

func HandleResponse(ctx *gin.Context, status int, response any) {
	ctx.JSON(status, response)
}

func Success(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, models.Response{
		Success: true,
		Status:  status,
		Data:    data,
	})
}

func Error(ctx *gin.Context, status int, message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	}

	ctx.JSON(status, models.Response{
		Success: false,
		Status:  status,
		Message: message,
	})
}

func AbortWithError(ctx *gin.Context, status int, message string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	}

	ctx.AbortWithStatusJSON(status, models.Response{
		Success: false,
		Status:  status,
		Message: message,
	})
}
