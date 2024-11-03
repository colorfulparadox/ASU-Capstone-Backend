package routes

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
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
	Name            bool `json:"name"`
	Username        bool `json:"username"`
	Password        bool `json:"password"`
	PermissionLevel bool `json:"permission_level"`
	Email           bool `json:"email"`
}

func Update_User(gc *gin.Context) {
	var updatedData UpdatedData
	var updateResult UpdateResult

	updatedData.PermissionLevel = -1

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&updatedData)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var update_error []error

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	if updatedData.Name != "" {
		update_error = append(update_error, database.Set_Name(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Name))
		if update_error[len(update_error)-1] == nil {
			updateResult.Name = true
		}
	}

	if updatedData.Username != "" {
		update_error = append(update_error, database.Set_Username(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Username))
		if update_error[len(update_error)-1] == nil {
			updateResult.Username = true
		}
	}

	if updatedData.Password != "" {
		update_error = append(update_error, database.Set_Password(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Password))
		if update_error[len(update_error)-1] == nil {
			updateResult.Password = true
		}
	}

	if updatedData.PermissionLevel >= 0 {
		update_error = append(update_error, database.Set_Permissions(updatedData.UserAuthID, updatedData.Edit_User, updatedData.PermissionLevel))
		if update_error[len(update_error)-1] == nil {
			updateResult.PermissionLevel = true
		}
	}

	if updatedData.Email != "" {
		update_error = append(update_error, database.Set_Email(updatedData.UserAuthID, updatedData.Edit_User, updatedData.Email))
		if update_error[len(update_error)-1] == nil {
			updateResult.Email = true
		}
	}

	for i := 0; i < len(update_error); i++ {
		if update_error[i] != nil {
			gc.Header("backend-error", update_error[i].Error())
		}
	}

	// Returns userResult
	gc.JSON(http.StatusOK, updateResult)
}
