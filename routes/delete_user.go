package routes

import (
	"BackEnd/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteUserTokens struct {
	CurrentAuthID string `json:"authID"`
	DeleteUser    string `json:"username"`
}

// Meant to be used in scenarios where you need to log out all devices
func Delete_User(gc *gin.Context) {
	var deleteUserTokens DeleteUserTokens

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&deleteUserTokens)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	user_creation_success := database.Delete_User(deleteUserTokens.CurrentAuthID, deleteUserTokens.DeleteUser)

	// Checks if there was an error
	switch user_creation_success {
	case 1:
		// This is an easter egg cause it should never get to this point because deleting should never
		// cause data mismatch and if there is something has gone terribly wrong
		log.Println("Something has gone terribly wrong")
	case 2:
		log.Println("Invalid Permissions")
	}

	// Puts int into JSON object
	userResult := UserResult{
		Result: user_creation_success,
	}

	UserResults(userResult.Result)

	// Returns userResult
	gc.JSON(http.StatusOK, userResult)
}
