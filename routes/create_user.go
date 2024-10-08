package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserData struct {
	AdminAuthID     string `json:"authID"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PermissionLevel int    `json:"permission_level"`
	Email           string `json:"email"`
}

// Creates users
func Create_User(gc *gin.Context) {
	var createUserData CreateUserData

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&createUserData)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	user_creation_success := database.New_User(createUserData.AdminAuthID, createUserData.Name, createUserData.Username, createUserData.Password, createUserData.PermissionLevel, createUserData.Email)

	// Puts int into JSON object
	userResult := UserResult{
		Result: user_creation_success,
	}

	// Checks if there was an error
	if !UserResults(userResult.Result) {
		gc.Request.Header.Add("backend-error", "true")
		gc.JSON(http.StatusForbidden, userResult)
	}

	// Returns userResult
	gc.JSON(http.StatusOK, userResult)
}
