package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingTest struct {
	Ping string `json:"ping"`
}

func Ping(gc *gin.Context) {
	pingTest := PingTest{
		Ping: "pong",
	}

	gc.JSON(http.StatusOK, pingTest)
}
