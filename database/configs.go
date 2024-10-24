package database

// An enum to show the result of a transaction
const (
	Success = iota
	Data_Already_Exists
	Incorrect_Permissions
	Invalid_Data
	Error_With_External_Service
	File_Error
)

// Config options for some things in the program
const (
	// Amount of time in half seconds to wait for response to be completed before timing out
	Completed_Timeout = 10
	// Model that will run the interactions, list of models available here: https://platform.openai.com/docs/models
	Model_Name = "gpt-4-turbo-preview"
	// Name to the Temp folder
	Temp_Path = "temp"
	// Name of the Menu folder
	Menu_Path = "menu"
	// Name of the Threads folder
	Threads_Path = "threads"
)
