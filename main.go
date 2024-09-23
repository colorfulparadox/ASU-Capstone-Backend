package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"BackEnd/database"
	"BackEnd/router"
	"BackEnd/routes"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// https://go.dev/doc/tutorial/web-service-gin

type message struct {
	Message string `json:"message"`
}

func main() {
	database.Randomize_auth_token("zB82b5xGHVoQiBz2NragQvV4Z5-Gy4aT5xJzCrOKtp0=")
	//database.New_User_From_Object(database.Verify_User_Login("johndoe", "securepassword"), user)
	//return
	databaseUrl := "postgres://project-persona:T%7D%3F_%5D0Lu8I98@postgres.blusnake.net:35432/project-persona"

	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	r := router.NewRouter("localhost", "4040", conn)

	router.AddRoute(&r, router.Receiver{
		Route:     "/getmsg",
		RouteType: router.RouteGet,
		Sender: func(gc *gin.Context, pool *pgxpool.Pool) {
			var msg message
			err := pool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&msg.Message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed %v\n", err)
			}
			gc.IndentedJSON(http.StatusOK, msg)
		},
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/login",
		RouteType: router.RoutePost,
		Sender:    routes.Login,
	})

	/*
		add_route(&router, Receiver{
			route:      "/login",
			routeType:  RoutePost,
			middleware: default_middleware,
			sender: func(gc *gin.Context, pool *pgxpool.Pool) {
				fmt.Println("login req")

				var loginReq LoginRequest

				err := gc.ShouldBindJSON(&loginReq)
				if err != nil {
					gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				fmt.Println(loginReq)

				gc.JSON(http.StatusOK, "{\"auth\":\"thisisakey123\"}")
			},
		})
	*/

	router.RunRouter(r)
}
