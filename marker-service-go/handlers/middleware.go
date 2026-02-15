package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	userIDHeader = "X-User-Id"
	userIDKey    = "userId"
)

// ExtractUserID middleware extracts and validates the X-User-Id header
func ExtractUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader(userIDHeader)
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("missing required header: %s", userIDHeader),
			})
			c.Abort()
			return
		}

		userId, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userId <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid %s header: must be a positive integer", userIDHeader),
			})
			c.Abort()
			return
		}

		// Store userId in context for handlers to use
		c.Set(userIDKey, userId)
		c.Next()
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(c *gin.Context) (int64, error) {
	userIDValue, exists := c.Get(userIDKey)
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	userId, ok := userIDValue.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type in context")
	}

	return userId, nil
}
