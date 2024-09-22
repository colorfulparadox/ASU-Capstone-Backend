package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginToken struct {
	AuthID     string `json:"authID"`
	DateIssued int64  `json:"dateIssued"`
	Expires    int64  `json:"expires"`
}

func Login(gc *gin.Context, pool *pgxpool.Pool) {
	var loginReq LoginRequest

	err := gc.ShouldBindJSON(&loginReq)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var username string = ""
	var password string

	pool.QueryRow(
		context.Background(),
		"SELECT username, password FROM users WHERE username = $1",
		loginReq.User,
	).Scan(&username, &password)

	fmt.Println("From the database:")
	fmt.Println(username)
	fmt.Println(password)

	if loginReq.Pass != password || username == "" {
		fmt.Println("invalid password")
		gc.JSON(http.StatusForbidden, "{}")
	}

	loginToken := LoginToken{
		AuthID:     "12asdasd3",
		DateIssued: time.Now().Unix(),
		Expires:    time.Now().Add(3 * 24 * time.Hour).Unix(),
	}

	gc.JSON(http.StatusOK, loginToken)
}
