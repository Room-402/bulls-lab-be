package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture User ID (from a header for now)
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		// Inject into context
		ctx := context.WithValue(c.Request.Context(), "audit_user_id", userID)
		ctx = context.WithValue(ctx, "audit_ip", c.ClientIP())

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
