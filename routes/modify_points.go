package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Sentiment = iota
	Time
	Knowledge
	Salesmanship
)

type SelectPoints struct {
	AdminAuthID      string `json:"authID"`
	Sentiment_Points int    `java:"sentiment_points"`
	Sales_Points     int    `java:"sales_points"`
	Knowledge_Points int    `java:"knowledge_points"`
}

type TotalPoints struct {
	Points int `json:"points"`
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
	points := database.Modify_Points(selectPoints.AdminAuthID, selectPoints.Sentiment_Points, selectPoints.Sales_Points, selectPoints.Knowledge_Points)

	if points < 0 {
		gc.Request.Header.Add("backend-error", "true")
		gc.JSON(http.StatusForbidden, "{}")
		return
	} else {
		totalPoints.Points = points
	}

	// Returns userResult
	gc.JSON(http.StatusOK, totalPoints)
}
