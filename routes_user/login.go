package routes_user

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

// Holds the auth_token
type LoginToken struct {
	AuthID string `json:"authID"`
}

func Login(gc *gin.Context) {
	var loginReq LoginRequest

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&loginReq)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the auth_token for the specifc user
	auth_token, err := database.Verify_User_Login(loginReq.User, loginReq.Pass)

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Puts auth_token into JSON object
	loginToken := LoginToken{
		AuthID: auth_token,
	}

	// Returns loginToken
	gc.JSON(http.StatusOK, loginToken)
}
