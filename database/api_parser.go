//This file is meant for interpreting the data from pstgres_functions.go and no direct connections to the server should be made through this file
//Any direct access of data should be handled by the postgres_functions.go file

package database

import (
	"time"
)

// An enum for the security level of certain actions
const (
	edit_self    = 0
	delete_self  = 0
	edit_users   = 1
	create_users = 1
	delete_users = 1
)

func Verify_User_Login(username string, password string) string {
	user := Retrieve_User_Username(username)
	if Retrieve_User_Username(username).Password == password {
		return user.AuthToken
	} else {
		return ""
	}
}

// Won't work until I rebuild the array cause rn auth keys arent' unique but verifies it's a real auth key and checks that it's not expired
func Verify_User_Auth_Token(auth_token string) bool {
	user := Retrieve_User_Auth_Token(auth_token)
	if user.AuthToken == auth_token {
		if user.DateExpr.Before(time.Now()) {
			return true
		}

		Randomize_Auth_Token_Auth_Token(auth_token)
	}

	return false
}

func Verify_Permissions(auth_token string, security_level int) bool {
	user := Get_User(auth_token)
	if user.PermissionLevel >= security_level {
		return true
	}

	return false
}

// Gets user data from an auth_token and verifies the
func Get_User(auth_token string) User {
	var user User
	if Verify_User_Auth_Token(auth_token) {
		return Retrieve_User_Auth_Token(auth_token)
	} else {
		return user
	}
}

func New_User(current_username, auth_token, name, username, password string, permission_level int, email string) bool {

	if Verify_Permissions(auth_token, create_users) {
		user := User{
			Name:            name,
			Username:        username,
			Password:        password,
			Points:          0,
			PermissionLevel: permission_level,
			Email:           email,
		}

		Create_User(user)

		return true
	}

	return false
}

func Set_Permissions(auth_token string, username string, permission int) {
	if Verify_Permissions(auth_token, edit_users) {
		user := Get_User(Retrieve_User_Username(username).AuthToken)
		user.PermissionLevel = permission
		Update_User(username, user)
	}
}
