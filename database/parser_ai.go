package database

// func New_AI(authID, name, instructions string) error {
// 	if !Verify_Permissions(authID, set_ai) {
// 		return Invalid_Permissions()
// 	}
// 	var ai AI
// 	ai.Name = name
// 	return create_ai(ai)

// }

// func Add_Menu(authID, ai_name, menu_data string) error {
// 	if !Verify_Permissions(authID, set_menu) {
// 		return Invalid_Permissions()
// 	}

// 	log.Println(menu_data)

// 	menu_name, err := create_menu_record(menu_data, ai_name+"_menu")
// 	if err != nil {
// 		log.Println(err)
// 		return Invalid_Data()
// 	}

// 	ai, err := retrieve_ai(ai_name)
// 	if err != nil {
// 		log.Println("Error retrieveing the ai:", err)
// 		return Invalid_Data()
// 	}

// 	err = upload_menu(menu_name, ai)
// 	if err != nil {
// 		log.Println("Error uploading the menu:", err)
// 		return Error_With_External_Service()
// 	}

// 	return nil
// }
