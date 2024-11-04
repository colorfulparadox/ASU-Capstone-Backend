package routes_persona

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Creates users
func End_Conversation(gc *gin.Context) {
	var conversation Conversation

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&conversation)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to results
	if conversation.ConversationID == "" {
		err = database.End_All_Persona_Conversations(conversation.AuthID)
	} else {
		err = database.End_Persona_Conversation(conversation.AuthID, conversation.ConversationID)
	}

	new_conversation := Conversation{
		ConversationID: conversation.ConversationID,
	}

	// Checks if there was an error
	if err != nil {
		gc.Header("backend-error", err.Error())
		gc.JSON(http.StatusForbidden, "{}")
		return
	}

	// Returns userResult
	gc.JSON(http.StatusOK, new_conversation)
}
