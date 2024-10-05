package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResetAuthTokens struct {
	CurrentAuthID string `json:"authID"`
	ResetUser     string `json:"username"`
}

// Meant to be used in scenarios where you need to log out all devices
func Reset_Auth_Token(gc *gin.Context) {
	var resetAuthTokens ResetAuthTokens

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&resetAuthTokens)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to creation (can be found in api_parser starting at line 23)
	user_creation_success := database.Randomize_Auth_Token(resetAuthTokens.CurrentAuthID, resetAuthTokens.ResetUser)

	// Puts int into JSON object
	userResult := UserResult{
		Result: user_creation_success,
	}

	// Checks if there was an error
	UserResults(user_creation_success)

	// Returns userResult
	gc.JSON(http.StatusOK, userResult)
}
