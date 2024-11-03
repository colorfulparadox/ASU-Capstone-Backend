package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Creates users
func Start_Conversation(gc *gin.Context) {
	var conversation Conversation

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&conversation)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results
	conversation.Message, err = database.Start_Persona_Conversation(conversation.AuthID, conversation.Message, conversation.ConversationID)

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, conversation)
}
