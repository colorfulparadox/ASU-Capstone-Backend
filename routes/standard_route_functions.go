// Basic file simply to hold commonly used data for other routes

package routes

import "log"

// StandardResult is a basic JSON format for returning one of the result enum types ()
type StandardResult struct {
	Result int `json:"result"`
}

type Conversation struct {
	AuthID         string `json:"authID"`
	Message        string `json:"message"`
	ConversationID string `json:"conversationID"`
}

func UserResults(user_creation_success int) bool {
	switch user_creation_success {
	case 0:
		log.Println("Successful")
		return true
	case 1:
		log.Println("User already exists")
		return false
	case 2:
		log.Println("Invalid Permissions")
		return false
	case 3:
		log.Println("Selected User Does Not Exist")
		return false
	}

	return false
}
