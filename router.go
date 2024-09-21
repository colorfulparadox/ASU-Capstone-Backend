package main

import (
	"fmt"

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

func new_router(
	ip string,
	port string,
	pool *pgxpool.Pool,
) Router {
	router := gin.Default()
	return Router{ip: ip, port: port, pool: pool, router: router}
}

func add_route(router *Router, receiver Receiver) {
	router.routes = append(router.routes, receiver)
}

// https://github.com/gin-gonic/examples/blob/master/cookie/main.go
func verifyCookies(router Router) gin.HandlerFunc {
	return func(gc *gin.Context) {
		gc.Next()
	}
}

func run_router(router Router) {
	fmt.Println(len(router.routes))
	for i := 0; i < len(router.routes); i++ {
		receiver := router.routes[i]

		if receiver.routeType == RouteGet {
			if receiver.authRequired {
				router.router.GET(receiver.route, verifyCookies(router), func(gc *gin.Context) {
					receiver.sender(gc, router.pool)
				})
				continue
			}
			//for no cookie requests
			router.router.GET(receiver.route, func(gc *gin.Context) {
				receiver.sender(gc, router.pool)
			})
		}

	}

	router.router.Run(router.ip + ":" + router.port)
}
