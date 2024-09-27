package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SelectPoints struct {
	AdminAuthID string `json:"authID"`
	Points      int    `java:"points"`
}

type TotalPoints struct {
	Verified bool `json:"verified"`
	Points   int  `json:"points"`
}

// Creates users
func ModifyPoints(gc *gin.Context) {
	var selectPoints SelectPoints
	var totalPoints TotalPoints

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&selectPoints)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	points := database.Modify_Points(selectPoints.AdminAuthID, selectPoints.Points)

	if points < 0 {
		totalPoints.Verified = false
	} else {
		totalPoints.Verified = true
		totalPoints.Points = points
	}

	// Returns userResult
	gc.JSON(http.StatusOK, totalPoints)
}
