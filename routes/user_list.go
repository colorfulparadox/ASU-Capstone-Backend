package routes

import (
	"BackEnd/database"
	"encoding/json"
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

// Creates users
func UserList(gc *gin.Context) {
	var specifyUser SpecifyUser

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&specifyUser)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Provided Json could not be parsed:", err)
		return
	}

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	if specifyUser.AdminList {
		users := database.Get_Admin_User_List(specifyUser.AuthID)
		var returnUserElement ReturnUserElement
		var userList []ReturnUserElement
		for i := 0; i < len(users); i++ {
			returnUserElement.Name = users[i][0]
			returnUserElement.Username = users[i][1]
			returnUserElement.Points, err = strconv.Atoi(users[i][2])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				return
			}
			userList = append(userList, returnUserElement)
		}

		jsonUserList, err := json.MarshalIndent(userList, "", "    ")
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Could not compile the user list:", err)
			return
		}

		// Returns userResult
		gc.JSON(http.StatusOK, jsonUserList)

	} else {
		users := database.Get_User_List(specifyUser.AuthID)
		var returnAdminUserElement ReturnAdminUserElement
		var userList []ReturnAdminUserElement
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

		jsonUserList, err := json.MarshalIndent(userList, "", "    ")
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("Could not compile the user list:", err)
			return
		}

		// Returns userResult
		gc.JSON(http.StatusOK, jsonUserList)
	}
}
