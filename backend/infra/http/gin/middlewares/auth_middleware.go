// infra/http/gin/middlewares/auth.go
package middlewares

import (
	// "fmt"
	"net/http"
	"strings"

	"ai_hub.com/app/core/ports/adminports"
	"github.com/gin-gonic/gin"
)

type TokenVerifier interface {
	Verify(token string) (*adminports.TokenPayload, error)
}

type UserPayload struct {
	ID    string   `json:"_id"`
	Email *string  `json:"email,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

func Auth(verifier TokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You don't have access"})
			return
		}

		token := strings.TrimSpace(authz)

		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = strings.TrimSpace(token[7:])
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You don't have access"})
			return
		}

		payload, err := verifier.Verify(token)
		if err != nil || payload == nil || payload.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Your access failed"})
			return
		}

		c.Set("userID", payload.UserID)

		c.Next()
	}
}
