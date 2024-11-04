package routes_user

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

	// Gets the enum int relating to results
	err = database.New_User(createUserData.AdminAuthID, createUserData.Name, createUserData.Username, createUserData.Password, createUserData.PermissionLevel, createUserData.Email)

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, StandardResult{Result: 0})
}
