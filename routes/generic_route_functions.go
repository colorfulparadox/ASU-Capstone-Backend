// Basic file simply to hold commonly used data for other routes

package routes

import "log"

// UserResult is a basic JSON format for returning one of the result enum types ()
type UserResult struct {
	Result int `json:"result"`
}

func UserResults(user_creation_success int) {
	switch user_creation_success {
	case 0:
		log.Println("Successful")
	case 1:
		log.Println("User already exists")
	case 2:
		log.Println("Invalid Permissions")
	case 3:
		log.Println("Selected User Does Not Exist")
	}
}
