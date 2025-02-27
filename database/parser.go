package database

import "os"

func Create_Tables() {
	create_users_table()
	// create_ai_table()
	// create_persona_table()
}

func Initalize_Directories() {
	os.MkdirAll(Conversation_Path, os.ModePerm)
	os.MkdirAll(Menu_Path, os.ModePerm)
}
