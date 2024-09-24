package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserData struct {
	AdminAuthID     string `json:"authID"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PermissionLevel int    `json:"permission_level"`
	Email           string `json:"email"`
}

// Creates users
func Create_User(gc *gin.Context) {
	var userData UserData

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&userData)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	user_creation_success := database.New_User(userData.AdminAuthID, userData.Name, userData.Username, userData.Password, userData.PermissionLevel, userData.Email)

	// Checks if there was an error
	UserResults(user_creation_success)

	// Puts int into JSON object
	userResult := UserResult{
		Result: user_creation_success,
	}

	// Returns userResult
	gc.JSON(http.StatusOK, userResult)
}
