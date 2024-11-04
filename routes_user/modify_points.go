package routes_user

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

type AddPoints struct {
	AdminAuthID      string `json:"authID"`
	Sentiment_Points int    `java:"sentiment_points"`
	Sales_Points     int    `java:"sales_points"`
	Knowledge_Points int    `java:"knowledge_points"`
}

type CurrentPoints struct {
	Sentiment_Points int `java:"sentiment_points"`
	Sales_Points     int `java:"sales_points"`
	Knowledge_Points int `java:"knowledge_points"`
}

// Creates users
func Modify_Points(gc *gin.Context) {
	var addPoints AddPoints
	var currentPoints CurrentPoints

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&addPoints)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	currentPoints.Sentiment_Points, currentPoints.Sales_Points, currentPoints.Knowledge_Points, err = database.Modify_Points(addPoints.AdminAuthID, addPoints.Sentiment_Points, addPoints.Sales_Points, addPoints.Knowledge_Points)

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, currentPoints)
}
