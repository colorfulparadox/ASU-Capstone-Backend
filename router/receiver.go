package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteType uint16

const (
	RouteGet RouteType = iota
	RoutePost
)

type senderFunc func(gc *gin.Context)

type middlewareFunc func(gc *gin.Context)

type Receiver struct {
	Route      string
	RouteType  RouteType
	Middleware func(gc *gin.Context)
	Sender     func(gc *gin.Context)
}

func default_middleware() middlewareFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func User_Role_Middleware(requiredRole string) middlewareFunc {
	return func(c *gin.Context) {
		// Assume we get the role from some context or header (example purpose)
		userRole := c.GetHeader("User")

		if userRole == "" || userRole != requiredRole {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}
