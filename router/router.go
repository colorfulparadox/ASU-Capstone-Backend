package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	port   string
	router *gin.Engine
}

func NewRouter(port string) Router {
	router := gin.New()

	// config := cors.DefaultConfig()
	// config.AllowOrigins = append(config.AllowOrigins, "*")
	// //config.AllowAllOrigins = true //we will want to disable that some day
	// config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// router.Use(cors.New(config))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	return Router{port: port, router: router}
}

func AddRoute(r *Router, receiver Receiver) {
	//router.routes = append(router.routes, receiver)

	if receiver.Middleware == nil {
		receiver.Middleware = default_middleware()
	}

	if receiver.RouteType == RoutePost {
		r.router.POST(receiver.Route, receiver.Middleware, receiver.Sender)
	} else {
		r.router.GET(receiver.Route, receiver.Middleware, receiver.Sender)
	}
}

func RunRouter(r Router) {
	r.router.Run(r.port)
}
