package routes_user

import (
	"BackEnd/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateResult struct {
	Name            bool `json:"name"`
	Username        bool `json:"username"`
	Password        bool `json:"password"`
	PermissionLevel bool `json:"permission_level"`
	Email           bool `json:"email"`
}

func Update_User(gc *gin.Context) {
	var userData UserData
	var updateResult UpdateResult

	userData.PermissionLevel = -1

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&userData)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var update_error []error

	// Gets the enum int relating to results (can be found in api_parser starting at line 23)
	if userData.Name != "" {
		update_error = append(update_error, database.Set_Name(userData.AuthID, userData.Edit_User, userData.Name))
		if update_error[len(update_error)-1] == nil {
			updateResult.Name = true
		}
	}

	if userData.Username != "" {
		update_error = append(update_error, database.Set_Username(userData.AuthID, userData.Edit_User, userData.Username))
		if update_error[len(update_error)-1] == nil {
			updateResult.Username = true
		}
	}

	if userData.Password != "" {
		update_error = append(update_error, database.Set_Password(userData.AuthID, userData.Edit_User, userData.Password))
		if update_error[len(update_error)-1] == nil {
			updateResult.Password = true
		}
	}

	if userData.PermissionLevel >= 0 {
		update_error = append(update_error, database.Set_Permissions(userData.AuthID, userData.Edit_User, userData.PermissionLevel))
		if update_error[len(update_error)-1] == nil {
			updateResult.PermissionLevel = true
		}
	}

	if userData.Email != "" {
		update_error = append(update_error, database.Set_Email(userData.AuthID, userData.Edit_User, userData.Email))
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
