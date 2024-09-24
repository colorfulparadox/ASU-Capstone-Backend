package routes

import (
	"BackEnd/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResetAuthTokens struct {
	CurrentAuthID string `json:"authID"`
	ResetUser     string `json:"username"`
}

// Meant to be used in scenarios where you need to log out all devices
func Reset_Auth_Token(gc *gin.Context, pool *pgxpool.Pool) {
	var resetAuthTokens ResetAuthTokens

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&resetAuthTokens)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to creation (can be found in api_parser starting at line 23)
	user_creation_success := database.Randomize_Auth_Token(resetAuthTokens.CurrentAuthID, resetAuthTokens.ResetUser)

	// Checks if there was an error
	switch user_creation_success {
	case 1:
		// This is an easter egg cause it should never get to this point because the random numbers should never
		// be the same and if they are there should be a more substantial error handler that should catch it
		log.Println("Something has gone terribly wrong")
	case 2:
		log.Println("Invalid Permissions")
	}

	// Puts int into JSON object
	userResult := UserResult{
		Result: user_creation_success,
	}

	// Returns userResult
	gc.JSON(http.StatusOK, userResult)
}
