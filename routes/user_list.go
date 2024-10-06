package routes

import (
	"BackEnd/database"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SpecifyUser struct {
	AuthID    string `json:"authID"`
	AdminList bool   `json:"admin"`
}

type UserElement struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Permissions int    `json:"permission_level"`
	Points      int    `json:"points"`
}

type UserElementList []UserElement

// Creates users
func UserList(gc *gin.Context) {
	var specifyUser SpecifyUser
	var userElement UserElement
	var userList UserElementList
	var users [][]string

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&specifyUser)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Provided Json could not be parsed:", err)
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	if specifyUser.AdminList {
		users = database.Get_Admin_User_List(specifyUser.AuthID)

		for i := 0; i < len(users); i++ {
			userElement.Name = users[i][0]
			userElement.Username = users[i][1]
			userElement.Email = users[i][2]
			userElement.Permissions, err = strconv.Atoi(users[i][3])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}
			userElement.Points, err = strconv.Atoi(users[i][4])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}
			userList = append(userList, userElement)
		}
	} else {
		users = database.Get_User_List(specifyUser.AuthID)

		for i := 0; i < len(users); i++ {
			userElement.Name = users[i][0]
			userElement.Username = users[i][1]
			userElement.Points, err = strconv.Atoi(users[i][2])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}

			userList = append(userList, userElement)
		}
	}

	log.Println("All users added to user list")

	// Returns userResult
	gc.JSON(http.StatusOK, userList)
}
