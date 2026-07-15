package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"meal-planner/internal/db"
	"meal-planner/internal/models"
)

// AuthMiddleware validates the API key in the request header
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

		// Look up the key in the database
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

		// Store the user in the context for use in handlers
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("device_id", user.DeviceID)

		c.Next()
	}
}

// RequireAdmin restricts access to admin users; must run after AuthMiddleware
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUserFromContext(c)
		if user == nil || !user.IsAdmin {
			c.JSON(403, models.ErrorResponse{
				Error:   "forbidden",
				Message: "Admin access required",
				Code:    "ADMIN_REQUIRED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ErrorHandler handles validation errors
func ErrorHandler(c *gin.Context, err error) {
	c.JSON(400, models.ErrorResponse{
		Error:   "invalid_request",
		Message: fmt.Sprintf("Request validation failed: %v", err),
		Code:    "VALIDATION_ERROR",
	})
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, models.SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse sends an error response
func ErrorResponseJSON(c *gin.Context, statusCode int, errCode, errMsg string) {
	c.JSON(statusCode, models.ErrorResponse{
		Error:   errCode,
		Message: errMsg,
		Code:    errCode,
	})
}

// GetUserFromContext extracts the user from the context
func GetUserFromContext(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*models.User)
}
