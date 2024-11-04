package routes_user

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

type UserDataList []UserData

// Creates users
func User_List(gc *gin.Context) {
	var specifyUser SpecifyUser
	var userData UserData
	var userList UserDataList
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
			userData.Name = users[i][0]
			userData.Username = users[i][1]
			userData.Email = users[i][2]
			userData.PermissionLevel, err = strconv.Atoi(users[i][3])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				gc.Request.Header.Add("backend-error", "true")
				gc.JSON(http.StatusForbidden, "{}")
				return
			}
			userData.Sentiment_Points, err = strconv.Atoi(users[i][4])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				gc.Request.Header.Add("backend-error", "true")
				gc.JSON(http.StatusForbidden, "{}")
				return
			}
			userData.Sales_Points, err = strconv.Atoi(users[i][5])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				gc.Request.Header.Add("backend-error", "true")
				gc.JSON(http.StatusForbidden, "{}")
				return
			}
			userData.Knowledge_Points, err = strconv.Atoi(users[i][6])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				gc.Request.Header.Add("backend-error", "true")
				gc.JSON(http.StatusForbidden, "{}")
				return
			}
			userList = append(userList, userData)
		}
	} else {
		users = database.Get_User_List(specifyUser.AuthID)

		for i := 0; i < len(users); i++ {
			userData.Name = users[i][0]
			userData.Username = users[i][1]
			userData.Average_Points, err = strconv.Atoi(users[i][2])
			if err != nil {
				log.Println("Incorrect data from API parser:", err)
				gc.Request.Header.Add("backend-error", "true")
				gc.JSON(http.StatusForbidden, "{}")
				return
			}

			userList = append(userList, userData)
		}
	}

	log.Println("All users added to user list")

	// Returns userResult
	gc.JSON(http.StatusOK, userList)
}
