package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SelectUser struct {
	AdminAuthID string `json:"authID"`
}

type ReturnUserData struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	Points          int    `json:"points"`
	PermissionLevel int    `json:"permission_level"`
	Email           string `json:"email"`
}

// Creates users
func Authenticate(gc *gin.Context) {
	var selectUser SelectUser
	var returnUserData ReturnUserData

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&selectUser)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	user := database.Verify_User_Auth_Token(selectUser.AdminAuthID)

	if user.Username != "" {
		returnUserData.Name = user.Name
		returnUserData.Username = user.Username
		returnUserData.Points = user.Sentiment_Points
		returnUserData.Points = user.Sales_Points
		returnUserData.Points = user.Knowledge_Points
		returnUserData.PermissionLevel = user.PermissionLevel
		returnUserData.Email = user.Email
	} else {
		gc.Request.Header.Add("backend-error", "true")
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, returnUserData)
}
