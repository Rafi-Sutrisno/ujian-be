package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")
		// fmt.Println("Origin:", origin)
		allowedOrigins := map[string]bool{
			"https://34.128.84.215": true, // Change to your actual frontend domain
			"http://localhost:3000":     true, // Allow local dev
		}

		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin) // Allow only specific origins
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, X-SafeExamBrowser-BrowserExamKeyHash, x-safeexambrowser-configkeyhash, x-safeexambrowser-requesthash")

		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
