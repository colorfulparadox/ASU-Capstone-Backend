package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func default_middleware() middlewareFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func User_Role_Middleware(requiredRole string) middlewareFunc {
	return func(c *gin.Context) {
		// Assume we get the role from some context or header (example purpose)
		userRole := c.GetHeader("user")

		if userRole == "" || userRole != requiredRole {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}
