package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type New_Menu struct {
	AI_Name string `json:"aiID"`
	Menu    string `json:"menu"`
}

// Creates users
func Add_Menu(gc *gin.Context) {
	var new_menu New_Menu

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&new_menu)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results
	err = database.Add_Menu(new_menu.AI_Name, new_menu.Menu)

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, StandardResult{Result: 0})
}
