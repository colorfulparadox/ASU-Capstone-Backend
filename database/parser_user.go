//This file is meant for interpreting the data from pstgres_functions.go and no direct connections to the server should be made through this file
//Any direct access of data should be handled by the postgres_functions.go file

package database

import (
	"log"
	"strconv"
	"time"
)

// An enum for the security level of certain actions
const (
	create_users        = 1
	get_users           = 1
	get_user_list       = 0
	get_admin_user_list = 1
	get_persona         = 0
	edit_persona        = 1
	get_ai              = 0
	edit_ai             = 1
	edit_self           = 0
	edit_users          = 1
	edit_permissions    = 1
	delete_self         = 0
	delete_users        = 1
)

// Verifies a login attempt with the username and password
func Verify_User_Login(username string, password string) string {
	user := retrieve_user_username(username)
	if VerifyPassword(user.PasswordHash, password) {
		log.Printf("Verified Login by: %s\n", username)
		return user.AuthToken
	} else {
		log.Println("Incorrect cridentials")
		return ""
	}
}

// Verifies an authentication attempt with the auth_token
// If successful returns the specified user object
// If unsuccessful returns an empty user object
func Verify_User_Auth_Token(auth_token string) User {
	user := retrieve_user_auth_token(auth_token)
	if user.AuthToken == auth_token {
		if user.DateExpr.After(time.Now().UTC()) {
			log.Printf("Verified Authentication of: %s\n", user.Username)
			log.Printf("Passing Current User to Function\n")
			return user
		}
		End_All_Persona_Conversations(auth_token)
		randomize_auth_token(auth_token)
		log.Println("Auth Token reset")
	}
	log.Println("Invalid Auth Token")
	var empty_user User
	empty_user.PermissionLevel = -1
	return empty_user
}

// Verifies the current user has the permissions to perform an action
func Verify_Permissions(auth_token string, security_level int) bool {
	user := Verify_User_Auth_Token(auth_token)

	if user.PermissionLevel == -1 {
		log.Println("Invalid permissions")
		return false
	}

	if user.PermissionLevel >= security_level {
		log.Println("Correct Permissions")
		return true
	} else {
		log.Printf("User: %s is not allowed to perform that action\n", user.Username)
		return false
	}
}

func Get_User_List(auth_token string) [][]string {
	var users [][]string
	if Verify_Permissions(auth_token, get_user_list) {
		userList := retrieve_user_list()
		for i := 0; i < len(userList); i++ {
			users = append(users, []string{userList[i].Name, userList[i].Username, strconv.Itoa((userList[i].Sentiment_Points + userList[i].Sales_Points + userList[i].Knowledge_Points) / 3)})
		}

		log.Println("All users returned")
	}

	return users
}

func Get_Admin_User_List(auth_token string) [][]string {
	var users [][]string
	if Verify_Permissions(auth_token, get_admin_user_list) {
		userList := retrieve_user_list()
		for i := 0; i < len(userList); i++ {
			users = append(users, []string{userList[i].Name, userList[i].Username, userList[i].Email, strconv.Itoa(userList[i].PermissionLevel), strconv.Itoa(userList[i].Sentiment_Points), strconv.Itoa(userList[i].Sales_Points), strconv.Itoa(userList[i].Knowledge_Points)})
		}
	}

	return users
}

// Adds a user to the database
func New_User(auth_token, name, username, password string, permission_level int, email string) int {

	if Verify_Permissions(auth_token, create_users) {
		user := User{
			Name:            name,
			Username:        username,
			Password:        password,
			PermissionLevel: permission_level,
			Email:           email,
		}

		if create_user(user) {
			return Success
		}
		return Data_Already_Exists
	}
	return Incorrect_Permissions
}

// Adds a user to the database from a user object
// This is used for testing and is not meant to be used in production
func New_User_From_Object(auth_token string, user User) bool {

	log.Printf("User permissions: %d\n", Verify_User_Auth_Token(auth_token).PermissionLevel)
	if Verify_Permissions(auth_token, create_users) {

		create_user(user)

		return true
	}

	return false
}

// Verifies auth token then returns user based on username if verified
func Get_User(auth_token, username string) User {
	var user User
	if Verify_Permissions(auth_token, get_users) {
		user = Verify_User_Auth_Token(auth_token)
		log.Printf("Passing Requested User to Function\n")
	}
	return user
}

func Verify_Request(auth_token, username string, user_permission, admin_permission int) (User, int) {
	//Determines if user is editing themselves or someone else and sets permissions accordingly
	var security_level int
	user := Verify_User_Auth_Token(auth_token)

	if user.PermissionLevel == -1 {
		return user, -1
	}

	if user.Username == username || username == "" {
		username = ""
		security_level = edit_self
	} else {
		security_level = edit_users
	}

	return user, security_level
}

// takes the current user's auth token and the username and new name of the user to be changed
func Set_Name(auth_token string, username string, new_name string) int {
	user, security_level := Verify_Request(auth_token, username, edit_self, edit_users)
	if security_level < 0 {
		return Invalid_Data
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		if username != "" {
			user = retrieve_user_username(username)
		}

		user.Name = new_name
		user.Password = ""

		if update_user(user.Username, user) {
			log.Println("User Name Updated")
			return Success
		}
	}

	return Incorrect_Permissions
}

// takes the current user's auth token and the old and new username of the user to be changed
func Set_Username(auth_token string, username string, new_username string) int {
	user, security_level := Verify_Request(auth_token, username, edit_self, edit_users)
	if security_level < 0 {
		return Invalid_Data
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		if username != "" {
			user = retrieve_user_username(username)
		}

		old_username := user.Username
		user.Username = new_username
		user.Password = ""

		if update_user(old_username, user) {
			log.Println("User Username Updated")
			return Success
		} else {
			return Data_Already_Exists
		}
	}

	return Incorrect_Permissions
}

// takes the current user's auth token and the username and new password of the user to be changed
func Set_Password(auth_token string, username string, new_password string) int {
	user, security_level := Verify_Request(auth_token, username, edit_self, edit_users)
	if security_level < 0 {
		return Invalid_Data
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		if username != "" {
			user = retrieve_user_username(username)
		}

		user.Password = new_password

		if update_user(user.Username, user) {
			log.Println("User Password Updated")
			return Success
		}
	}

	return Incorrect_Permissions
}

// takes the current user's auth token and the username and new permission level of the user to be changed
func Set_Permissions(auth_token string, username string, new_permission int) int {
	//Checks permissions then applies change to user
	if Verify_Permissions(auth_token, edit_permissions) {
		user := retrieve_user_username(username)

		user.PermissionLevel = new_permission
		user.Password = ""

		if update_user(user.Username, user) {
			log.Println("User Permissions Updated")
			return Success
		} else {
			return Data_Already_Exists
		}
	}

	return Incorrect_Permissions
}

// takes the current user's auth token and the username and new email of the user to be changed
func Set_Email(auth_token string, username string, new_email string) int {
	user, security_level := Verify_Request(auth_token, username, edit_self, edit_users)
	if security_level < 0 {
		return Invalid_Data
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		if username != "" {
			user = retrieve_user_username(username)
		}

		user.Email = new_email
		user.Password = ""

		if update_user(user.Username, user) {
			log.Println("User Email Updated")
			return Success
		} else {
			return Data_Already_Exists
		}
	}

	return Incorrect_Permissions
}

func Modify_Points(auth_token string, sentiment_points, sales_points, knowledge_points int) int {
	//Validates token then adds points to user
	user := Verify_User_Auth_Token(auth_token)
	if Verify_Permissions(auth_token, edit_self) {

		user.Sentiment_Points += sentiment_points

		// Minimum number of points is 0
		if user.Sentiment_Points < 0 {
			user.Sentiment_Points = 0
		}

		user.Sales_Points += sales_points

		// Minimum number of points is 0
		if user.Sales_Points < 0 {
			user.Sales_Points = 0
		}

		user.Knowledge_Points += knowledge_points

		// Minimum number of points is 0
		if user.Knowledge_Points < 0 {
			user.Knowledge_Points = 0
		}

		if update_user(user.Username, user) {
			return Success
		} else {
			return Invalid_Data
		}
	}

	return Incorrect_Permissions
}

// Randomizes the user auth token
func Randomize_Auth_Token(auth_token, username string) int {
	user, security_level := Verify_Request(auth_token, username, edit_self, edit_users)
	if security_level < 0 {
		return Invalid_Data
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		if username != "" {
			user = retrieve_user_username(username)
		}

		randomize_auth_token(user.AuthToken)
		log.Println("Authentication Token Randomized")
		return Success
	}

	return Incorrect_Permissions
}

// Deletes a specified user
func Delete_User(auth_token, username string) int {
	user, security_level := Verify_Request(auth_token, username, delete_self, delete_users)
	if security_level < 0 {
		return Invalid_Data
	}

	//Applies change to user
	if Verify_Permissions(auth_token, security_level) {
		if username != "" {
			user = retrieve_user_username(username)
		}

		delete_user(user.Username)
		log.Println("User Deleted")
		return Success
	}

	return Incorrect_Permissions
}
