package interceptors

import "github.com/gin-gonic/gin"

const AppVersion = "1.0.0"

func VersionInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-API-Version", AppVersion)
		c.Next()
	}
}
