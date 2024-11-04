package database

import (
	"errors"
)

// Functions to show the result of a transaction
func Data_Already_Exists() error {
	return errors.New("data_already_exists")
}

func Invalid_Permissions() error {
	return errors.New("invalid_permissions")
}

func Invalid_Data() error {
	return errors.New("invalid_data")
}

func Error_With_External_Service() error {
	return errors.New("error_with_external_service")
}

func File_Error() error {
	return errors.New("file_error")
}

const (
	user  = 0
	admin = 1
)

// An enum for the security level of certain actions
// 0=
const (
	get_self            = user
	set_self            = user
	delete_self         = user
	get_users           = admin
	set_users           = admin
	delete_users        = admin
	get_user_list       = user
	get_admin_user_list = admin
	get_persona         = user
	set_persona         = admin
	get_ai              = user
	set_ai              = admin
	set_permissions     = admin
	set_conversation    = user
	set_menu            = admin
)

// Config options for some things in the program
const (
	// Amount of time in half seconds to wait for response to be completed before timing out
	Completed_Timeout = 20
	// Model that will run the interactions, list of models available here: https://platform.openai.com/docs/models
	Model_Name = "gpt-4-turbo-preview"
	// Name to the Temp folder
	Temp_Path = "data/temp"
	// Name of the Menu folder
	Menu_Path = "data/menu"
	// Name of the File ID folder
	FileID_Path = "data/file_ids"
	// Name of the Threads folder
	Conversation_Path = "data/conversations"
)
