//This file is meant for interpreting the data from pstgres_functions.go and no direct connections to the server should be made through this file
//Any direct access of data should be handled by the postgres_functions.go file

package database

import (
	"errors"
	"log"
	"strconv"
	"time"
)

// Verifies a login attempt with the username and password
func Verify_User_Login(username string, password string) (string, error) {
	user, err := retrieve_user_username(username)
	if err != nil {
		return "", err
	}
	if VerifyPassword(user.PasswordHash, password) {
		user, err = Verify_User_Auth_Token(user.AuthToken)
		if err != nil {
			user, _ = retrieve_user_username(username)
		}
		log.Printf("Verified Login by: %s\n", username)
		return user.AuthToken, err
	} else {
		log.Println("Incorrect cridentials")
		return "", nil
	}
}

// Verifies an authentication attempt with the auth_token
// If successful returns the specified user object
// If unsuccessful returns an empty user object
func Verify_User_Auth_Token(auth_token string) (User, error) {
	var empty_user User
	user, err := retrieve_user_auth_token(auth_token)
	if err != nil {
		return empty_user, err
	}
	if user.AuthToken == auth_token {
		if user.DateExpr.After(time.Now().UTC()) && Authentication_Token_Forced_Time_Reset != 0 {
			log.Printf("Verified Authentication of: %s\n", user.Username)
			log.Printf("Passing Current User to Function\n")
			return user, nil
		}
		randomize_auth_token(auth_token)
		log.Println("Auth Token reset")
	}
	log.Println("Invalid Auth Token")
	return empty_user, err
}

// Verifies the current user has the permissions to perform an action
func Verify_Permissions(auth_token string, security_level int) bool {
	user, err := Verify_User_Auth_Token(auth_token)
	if err != nil {
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
		userList, err := retrieve_user_list()
		if err != nil {
			return nil
		}
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
		userList, err := retrieve_user_list()
		if err != nil {
			return nil
		}
		for i := 0; i < len(userList); i++ {
			users = append(users, []string{userList[i].Name, userList[i].Username, userList[i].Email, strconv.Itoa(userList[i].PermissionLevel), strconv.Itoa(userList[i].Sentiment_Points), strconv.Itoa(userList[i].Sales_Points), strconv.Itoa(userList[i].Knowledge_Points)})
		}
	}

	return users
}

// Adds a user to the database
func New_User(auth_token, name, username, password string, permission_level int, email string) error {

	if Verify_Permissions(auth_token, set_users) {
		user := User{
			Name:            name,
			Username:        username,
			Password:        password,
			PermissionLevel: permission_level,
			Email:           email,
		}

		if create_user(user) == nil {
			return nil
		}
		return Data_Already_Exists()
	}
	return Invalid_Permissions()
}

// Adds a user to the database from a user object
// This is used for testing and is not meant to be used in production
func New_User_From_Object(auth_token string, user User) error {
	if Verify_Permissions(auth_token, set_users) {

		create_user(user)

		return nil
	}

	return Invalid_Data()
}

func Verify_Request(auth_token, username string, user_permission, admin_permission int) (User, error) {
	//Determines if user is editing themselves or someone else and sets permissions accordingly
	var security_level int
	user, err := Verify_User_Auth_Token(auth_token)
	if err != nil {
		return user, err
	}

	if user.Username == username || username == "" {
		username = ""
		security_level = set_self
	} else {
		security_level = set_users
	}

	if user.PermissionLevel < security_level {
		return user, Invalid_Permissions()
	}

	return user, nil
}

// takes the current user's auth token and the username and new name of the user to be changed
func Set_Name(auth_token string, username string, new_name string) error {
	user, err := Verify_Request(auth_token, username, set_self, set_users)
	if err != nil {
		return Invalid_Permissions()
	}

	//Applies change to user
	if username != "" {
		user, err = retrieve_user_username(username)
	}
	if err != nil {
		return Invalid_Data()
	}

	user.Name = new_name
	user.Password = ""

	if update_user(user.Username, user) == nil {
		log.Println("User Name Updated")
		return nil
	}

	return errors.New("User could not be updated")
}

// takes the current user's auth token and the old and new username of the user to be changed
func Set_Username(auth_token string, username string, new_username string) error {
	user, err := Verify_Request(auth_token, username, set_self, set_users)
	if err != nil {
		return Invalid_Permissions()
	}

	End_All_Persona_Conversations(username)

	//Applies change to user
	if username != "" {
		user, err = retrieve_user_username(username)
	}
	if err != nil {
		return Invalid_Data()
	}

	old_username := user.Username
	user.Username = new_username
	user.Password = ""

	if update_user(old_username, user) == nil {
		log.Println("User Username Updated")
		return nil
	} else {
		return Data_Already_Exists()
	}
}

// takes the current user's auth token and the username and new password of the user to be changed
func Set_Password(auth_token string, username string, new_password string) error {
	user, err := Verify_Request(auth_token, username, set_self, set_users)
	if err != nil {
		return Invalid_Permissions()
	}

	//Applies change to user
	if username != "" {
		user, err = retrieve_user_username(username)
	}
	if err != nil {
		return Invalid_Data()
	}

	user.Password = new_password

	if update_user(user.Username, user) == nil {
		log.Println("User Password Updated")
		return nil
	}

	return errors.New("User could not be updated")

}

// takes the current user's auth token and the username and new permission level of the user to be changed
func Set_Permissions(auth_token string, username string, new_permission int) error {
	//Checks permissions then applies change to user
	if Verify_Permissions(auth_token, set_permissions) {
		user, err := retrieve_user_username(username)
		if err != nil {
			return Invalid_Data()
		}

		user.PermissionLevel = new_permission
		user.Password = ""

		if update_user(user.Username, user) == nil {
			log.Println("User Permissions Updated")
			return nil
		} else {
			return Data_Already_Exists()
		}
	}
	return Invalid_Permissions()
}

// takes the current user's auth token and the username and new email of the user to be changed
func Set_Email(auth_token string, username string, new_email string) error {
	user, err := Verify_Request(auth_token, username, set_self, set_users)
	if err != nil {
		return Invalid_Permissions()
	}

	//Applies change to user
	if username != "" {
		user, err = retrieve_user_username(username)
	}

	if err != nil {
		return Invalid_Data()
	}

	user.Email = new_email
	user.Password = ""

	if update_user(user.Username, user) == nil {
		log.Println("User Email Updated")
		return nil
	} else {
		return Data_Already_Exists()
	}
}

func Modify_Points(auth_token string, sentiment_points, sales_points, knowledge_points int) (int, int, int, error) {
	//Validates token then adds points to user
	user, err := Verify_User_Auth_Token(auth_token)
	if err != nil {
		return 0, 0, 0, Invalid_Data()
	}
	if Verify_Permissions(auth_token, set_self) {

		user.Sentiment_Points = (user.Sentiment_Points + sentiment_points) / 2

		// Minimum number of points is 0
		if user.Sentiment_Points < 0 {
			user.Sentiment_Points = 0
		}

		user.Sales_Points = (user.Sales_Points + sales_points) / 2

		// Minimum number of points is 0
		if user.Sales_Points < 0 {
			user.Sales_Points = 0
		}

		user.Knowledge_Points = (user.Knowledge_Points + knowledge_points) / 2

		// Minimum number of points is 0
		if user.Knowledge_Points < 0 {
			user.Knowledge_Points = 0
		}

		if update_user(user.Username, user) == nil {
			return user.Sentiment_Points, user.Sales_Points, user.Knowledge_Points, nil
		} else {
			return 0, 0, 0, Invalid_Data()
		}
	}

	return 0, 0, 0, Invalid_Permissions()

}

// Randomizes the user auth token
func Randomize_Auth_Token(auth_token, username string) error {
	user, err := Verify_Request(auth_token, username, set_self, set_users)
	if err != nil {
		return Invalid_Permissions()
	}

	//Applies change to user
	if username != "" {
		user, err = retrieve_user_username(username)
	}

	if err != nil {
		return Invalid_Data()
	}

	randomize_auth_token(user.AuthToken)
	log.Println("Authentication Token Randomized")
	return nil
}

// Deletes a specified user
func Delete_User(auth_token, username string) error {
	user, err := Verify_Request(auth_token, username, delete_self, delete_users)
	if err != nil {
		return Invalid_Permissions()
	}

	//Applies change to user
	if username != "" {
		user, err = retrieve_user_username(username)
	}

	if err != nil {
		return Invalid_Data()
	}

	delete_user(user.Username)
	log.Println("User Deleted")
	return nil
}
