package middlewares

import (
	"log"
	"net/http"

	"ai_hub.com/app/infra/http/httperrormapper"
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
		case httperrormapper.IsAdminError(err):
			status, body = httperrormapper.MapAdminErrorToHttp(err)
		case httperrormapper.IsApiKeyError(err):
			status, body = httperrormapper.MapApiKeyErrorToHttp(err)
		case httperrormapper.IsProjectError(err):
			status, body = httperrormapper.MapProjectErrorToHttp(err)
		case httperrormapper.IsPromptError(err):
			status, body = httperrormapper.MapPromptErrorToHttp(err)
		case httperrormapper.IsTaskError(err):
			status, body = httperrormapper.MapTaskErrorToHttp(err)
		default:
			status, body = httperrormapper.FallbackInternal(err)
		}

		log.Printf("[http][error] %T: %v => %d", err, err, status)

		c.AbortWithStatusJSON(status, body)
	}
}
