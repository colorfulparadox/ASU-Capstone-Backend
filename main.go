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

type user struct {
	Username string `json:"id"`
	Points   int32  `json:"points"`
}

var users = []user{
	{Username: "Bob", Points: 5},
	{Username: "Jimmy", Points: 20},
	{Username: "Saul", Points: 10},
}

var dbPool *pgxpool.Pool

func main() {

	databaseUrl := "postgres://project-persona:T%7D%3F_%5D0Lu8I98@postgres.blusnake.net:35432/project-persona"

	// this returns connection pool
	//dbPool, err := pgx.Connect(context.Background(), databaseUrl)

	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	dbPool = conn

	/*
		var greeting string
		err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(greeting)
	*/

	router := gin.Default()
	router.GET("/getusers", getUsers)
	router.GET("/getmsg", getMessage)
	router.Run("localhost:8081")
}

func getUsers(gc *gin.Context) {
	gc.IndentedJSON(http.StatusOK, users)
	//context.JSON(http.StatusOK, users)
}

func getMessage(gc *gin.Context) {
	var msg message
	err := dbPool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&msg.Message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed %v\n", err)
	}

	gc.IndentedJSON(http.StatusOK, msg)
}
