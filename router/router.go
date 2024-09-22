package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	ip     string
	port   string
	routes []Receiver
	pool   *pgxpool.Pool
	router *gin.Engine
}

func NewRouter(
	ip string,
	port string,
	pool *pgxpool.Pool,
) Router {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = append(config.AllowOrigins, "*")
	//config.AllowAllOrigins = true //we will want to disable that some day
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	return Router{ip: ip, port: port, pool: pool, router: router}
}

func AddRoute(router *Router, receiver Receiver) {
	//router.routes = append(router.routes, receiver)

	if receiver.Middleware == nil {
		receiver.Middleware = default_middleware
	}

	if receiver.RouteType == RoutePost {
		router.router.POST(
			receiver.Route,
			func(gc *gin.Context) {
				receiver.Middleware(gc, router.pool)
			},
			func(gc *gin.Context) {
				receiver.Sender(gc, router.pool)
			},
		)
	} else {
		router.router.GET(
			receiver.Route,
			func(gc *gin.Context) {
				receiver.Middleware(gc, router.pool)
			},
			func(gc *gin.Context) {
				receiver.Sender(gc, router.pool)
			},
		)
	}
}

func RunRouter(router Router) {
	router.router.Run(router.ip + ":" + router.port)
}
