package routes

import (
	"BackEnd/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserData struct {
	AdminAuthID     string `json:"authID"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PermissionLevel int    `json:"permission_level"`
	Email           string `json:"email"`
}

type UserResult struct {
	Result int `json:"result"`
}

func Create_User(gc *gin.Context, pool *pgxpool.Pool) {
	var userData UserData

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&userData)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the enum int relating to creation (can be found in api_parser starting at line 23)
	user_creation_success := database.New_User(userData.AdminAuthID, userData.Name, userData.Username, userData.Password, userData.PermissionLevel, userData.Email)

	// Checks if user is valid
	if user_creation_success > 0 {
		if user_creation_success == 1 {
			log.Println("User already exists")
		}

		log.Println("Invalid Permissions")
	}

	// Puts int into JSON object
	userResult := UserResult{
		Result: user_creation_success,
	}

	// Returns userResult
	gc.JSON(http.StatusOK, userResult)
}
