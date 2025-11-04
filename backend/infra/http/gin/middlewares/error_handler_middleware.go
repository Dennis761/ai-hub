package middlewares

import (
	"log"
	"net/http"

	"ai_hub.com/app/infra/http/httpErrorMapper"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if err == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "unknown error",
				"key":   "InternalServerError",
			})
			return
		}

		var status int
		var body map[string]string

		switch {
		case httpErrorMapper.IsAdminError(err):
			status, body = httpErrorMapper.MapAdminErrorToHttp(err)
		case httpErrorMapper.IsApiKeyError(err):
			status, body = httpErrorMapper.MapApiKeyErrorToHttp(err)
		case httpErrorMapper.IsProjectError(err):
			status, body = httpErrorMapper.MapProjectErrorToHttp(err)
		case httpErrorMapper.IsPromptError(err):
			status, body = httpErrorMapper.MapPromptErrorToHttp(err)
		case httpErrorMapper.IsTaskError(err):
			status, body = httpErrorMapper.MapTaskErrorToHttp(err)
		default:
			status, body = httpErrorMapper.FallbackInternal(err)
		}

		log.Printf("[http][error] %T: %v => %d", err, err, status)

		c.AbortWithStatusJSON(status, body)
	}
}
