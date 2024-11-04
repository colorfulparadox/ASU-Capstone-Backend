package routes_user

import (
	"BackEnd/database"
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
	err = database.Delete_User(deleteUserTokens.CurrentAuthID, deleteUserTokens.DeleteUser)

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, StandardResult{Result: 0})
}
