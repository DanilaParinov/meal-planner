package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"meal-planner/internal/db"
	"meal-planner/internal/models"
)

// AuthMiddleware проверяет API ключ в заголовке
func AuthMiddleware(repo *db.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "X-API-Key header is required",
				Code:    "MISSING_API_KEY",
			})
			c.Abort()
			return
		}

		// Проверяем ключ в БД
		user, err := repo.GetUserByAPIKey(apiKey)
		if err != nil {
			c.JSON(401, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid API key",
				Code:    "INVALID_API_KEY",
			})
			c.Abort()
			return
		}

		// Сохраняем пользователя в контексте для использования в handlers'ах
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("device_id", user.DeviceID)

		c.Next()
	}
}

// ErrorHandler обрабатывает ошибки валидации
func ErrorHandler(c *gin.Context, err error) {
	c.JSON(400, models.ErrorResponse{
		Error:   "invalid_request",
		Message: fmt.Sprintf("Request validation failed: %v", err),
		Code:    "VALIDATION_ERROR",
	})
}

// SuccessResponse отправляет успешный ответ
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, models.SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse отправляет ошибку
func ErrorResponseJSON(c *gin.Context, statusCode int, errCode, errMsg string) {
	c.JSON(statusCode, models.ErrorResponse{
		Error:   errCode,
		Message: errMsg,
		Code:    errCode,
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*models.User)
}
