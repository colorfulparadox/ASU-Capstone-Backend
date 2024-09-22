package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// https://go.dev/doc/tutorial/web-service-gin

type message struct {
	Message string `json:"message"`
}

func main() {
	databaseUrl := "postgres://project-persona:T%7D%3F_%5D0Lu8I98@postgres.blusnake.net:35432/project-persona"

	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	router := new_router("192.168.0.60", "4040", conn)

	add_route(&router, Receiver{
		route:        "/getmsg",
		authRequired: false,
		routeType:    RouteGet,
		sender: func(gc *gin.Context, pool *pgxpool.Pool) {
			var msg message
			err := pool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&msg.Message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed %v\n", err)
			}
			gc.IndentedJSON(http.StatusOK, msg)
		},
	})

	add_route(&router, Receiver{
		route:        "/login",
		authRequired: false,
		routeType:    RoutePost,
		sender: func(gc *gin.Context, pool *pgxpool.Pool) {
			fmt.Println("login req")
			body := gc.Request.Body
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
			}
			fmt.Println(body)

			gc.JSON(http.StatusOK, "{\"auth\":\"thisisakey123\"}")
		},
	})

	run_router(router)
}
