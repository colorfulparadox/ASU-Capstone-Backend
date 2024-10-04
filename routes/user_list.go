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

type ReturnUserElement struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Points   int    `json:"points"`
}

type ReturnAdminUserElement struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Permissions int    `json:"permission_level"`
	Points      int    `json:"points"`
}

type ReturnUserList []ReturnUserElement

type ReturnAdminUserList []ReturnAdminUserElement

// Creates users
func UserList(gc *gin.Context) {
	var specifyUser SpecifyUser
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
		var returnAdminUserElement ReturnAdminUserElement
		var userList ReturnAdminUserList
		for i := 0; i < len(users); i++ {
			returnAdminUserElement.Name = users[i][0]
			returnAdminUserElement.Username = users[i][1]
			returnAdminUserElement.Email = users[i][2]
			returnAdminUserElement.Permissions, err = strconv.Atoi(users[i][3])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}
			returnAdminUserElement.Points, err = strconv.Atoi(users[i][4])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}
			userList = append(userList, returnAdminUserElement)
		}

		log.Println("All users added to user list")

		// Returns userResult
		gc.JSON(http.StatusOK, userList)
	} else {
		users = database.Get_User_List(specifyUser.AuthID)
		log.Println("Users added to users object")
		var returnUserElement ReturnUserElement
		var userList ReturnUserList

		for i := 0; i < len(users); i++ {
			log.Println("users name:", users[i][0])
			returnUserElement.Name = users[i][0]
			returnUserElement.Username = users[i][1]
			returnUserElement.Points, err = strconv.Atoi(users[i][2])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}

			userList = append(userList, returnUserElement)
		}

		log.Println("All users added to user list")

		// Returns userResult
		gc.JSON(http.StatusOK, userList)
	}
}
