//This file is meant for interpreting the data from pstgres_functions.go and no direct connections to the server should be made through this file
//Any direct access of data should be handled by the postgres_functions.go file

package database

import (
	"log"
	"time"
)

// An enum for the security level of certain actions
const (
	create_users     = 1
	get_users        = 1
	edit_self        = 0
	edit_users       = 1
	edit_permissions = 1
	delete_self      = 0
	delete_users     = 1
)

func Create_Tables() {
	create_users_table()
}

// Verifies a login attempt with the username and password
func Verify_User_Login(username string, password string) string {
	user := retrieve_user_username(username)
	if VerifyPassword(user.PasswordHash, password) {
		log.Println("Authenticated")
		return user.AuthToken
	} else {
		log.Println("Incorrect cridentials")
		return ""
	}
}

// Verifies a login attempt with the auth_token
func Verify_User_Auth_Token(auth_token string) bool {
	user := retrieve_user_auth_token(auth_token)
	if user.AuthToken == auth_token {
		log.Println("Expire date: ", user.DateExpr)
		log.Println("Current date: ", time.Now().UTC())
		if user.DateExpr.After(time.Now().UTC()) {
			log.Printf("Verified")
			return true
		}
		randomize_auth_token(auth_token)
	}

	return false
}

// Verifies the current user has the permissions to perform an action
func Verify_Permissions(auth_token string, security_level int) bool {
	user := Get_Self(auth_token)
	if user.PermissionLevel >= security_level {
		log.Println("Correct Permissions")
		return true
	} else {
		log.Println("Invalid Permissions")
		return false
	}
}

// Adds a user to the database
func New_User(auth_token, name, username, password string, permission_level int, email string) bool {

	if Verify_Permissions(auth_token, create_users) {
		user := User{
			Name:            name,
			Username:        username,
			Password:        password,
			Points:          0,
			PermissionLevel: permission_level,
			Email:           email,
		}

		create_user(user)

		return true
	}

	return false
}

// Adds a user to the database from a user object
func New_User_From_Object(auth_token string, user User) bool {

	log.Printf("User permissions: %d\n", retrieve_user_auth_token(auth_token).PermissionLevel)
	if Verify_Permissions(auth_token, create_users) {

		create_user(user)

		return true
	}

	return false
}

// Verifies auth token then returns self if Verified
func Get_Self(auth_token string) User {
	var user User
	if Verify_User_Auth_Token(auth_token) {
		return retrieve_user_auth_token(auth_token)
	} else {
		return user
	}
}

// Verifies auth token then returns user based on username if verified
func Get_User(auth_token, username string) User {
	var user User
	if Verify_Permissions(auth_token, get_users) {
		return retrieve_user_auth_token(auth_token)
	}
	return user
}

// takes the current user's auth token and the username and new name of the user to be changed
func Set_Name(auth_token string, username string, new_name string) {
	//Determines if user is editing themselves or someone else and sets permissions accordingly
	var security_level int
	if Get_Self(auth_token).Username == username {
		security_level = edit_self
	} else {
		security_level = edit_users
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		user := retrieve_user_username(username)
		user.Name = new_name
		if update_user(username, user) {
			log.Println("Updated")
		} else {
			log.Println("Not Updated")
		}
	}
}

// takes the current user's auth token and the old and new username of the user to be changed
func Set_Username(auth_token string, username string, new_username string) {
	//Determines if user is editing themselves or someone else and sets permissions accordingly
	var security_level int
	if Get_Self(auth_token).Username == username {
		security_level = edit_self
	} else {
		security_level = edit_users
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		user := retrieve_user_username(username)
		user.Username = new_username
		if update_user(username, user) {
			log.Println("Updated")
		} else {
			log.Println("Not Updated")
		}
	}
}

// takes the current user's auth token and the username and new password of the user to be changed
func Set_Password(auth_token string, username string, new_password string) {
	//Determines if user is editing themselves or someone else and sets permissions accordingly
	var security_level int
	if Get_Self(auth_token).Username == username {
		security_level = edit_self
	} else {
		security_level = edit_users
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		user := retrieve_user_username(username)
		user.Password = new_password
		if update_user(username, user) {
			log.Println("Updated")
		} else {
			log.Println("Not Updated")
		}
	}
}

// takes the current user's auth token and the username and new permission level of the user to be changed
func Set_Permissions(auth_token string, username string, new_permission int) {
	//Checks permissions then applies change to user
	if Verify_Permissions(auth_token, edit_permissions) {
		user := retrieve_user_username(username)
		user.PermissionLevel = new_permission
		if update_user(username, user) {
			log.Println("Updated")
		} else {
			log.Println("Not Updated")
		}
	}
}

// takes the current user's auth token and the username and new email of the user to be changed
func Set_Email(auth_token string, username string, new_email string) {
	//Determines if user is editing themselves or someone else and sets permissions accordingly
	var security_level int
	if Get_Self(auth_token).Username == username {
		security_level = edit_self
	} else {
		security_level = edit_users
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		user := retrieve_user_username(username)
		user.Email = new_email
		if update_user(username, user) {
			log.Println("Updated")
		} else {
			log.Println("Not Updated")
		}
	}
}
