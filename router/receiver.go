package router

import (
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
