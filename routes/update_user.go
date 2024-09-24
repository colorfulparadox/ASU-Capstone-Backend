package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdatedData struct {
	UserAuthID      string `json:"authID"`
	Edit_User       string `json:"edit_user"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PermissionLevel int    `json:"permission_level"`
	Email           string `json:"email"`
}

type UpdateResult struct {
	Result []int `json:"result"`
}

func Update_User(gc *gin.Context, pool *pgxpool.Pool) {
	var updatedData UpdatedData

	updatedData.PermissionLevel = -1

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&updatedData)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user_update_success []int

	// Gets the enum int relating to creation (can be found in api_parser starting at line 23)
	if updatedData.Name != "" {
		user_update_success = append(user_update_success, database.Set_Name(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Name))
	}

	if updatedData.Username != "" {
		user_update_success = append(user_update_success, database.Set_Username(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Username))
	}

	if updatedData.Password != "" {
		user_update_success = append(user_update_success, database.Set_Password(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Password))
	}

	if updatedData.PermissionLevel >= 0 {
		user_update_success = append(user_update_success, database.Set_Permissions(updatedData.UserAuthID, updatedData.Edit_User, updatedData.PermissionLevel))
	}

	if updatedData.Email != "" {
		user_update_success = append(user_update_success, database.Set_Email(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Email))
	}

	// Puts int into JSON object
	updateResult := UpdateResult{
		Result: user_update_success,
	}

	// Returns userResult
	gc.JSON(http.StatusOK, updateResult)
}
