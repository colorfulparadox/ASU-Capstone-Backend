package database

import (
	"log"
)

type Conversation struct {
	ConversationID string `json:"conversation_id"`
	AssistantID    string `json:"assistant_id"`
	ThreadID       string `json:"thread_id"`
	RunID          string `json:"run_id"`
}

// func New_Persona(persona_name, ai_name, description, instructions string) error {
// 	persona := Persona{
// 		Name:         persona_name,
// 		AIName:       ai_name,
// 		Description:  description,
// 		Instructions: instructions,
// 	}

// 	return create_persona(persona)
// }

// TODO: Finish function and possibly add ability to select specific persona model
func Start_Persona_Conversation(authID, instructions, message, conversation_id string) (string, error) {
	if !Verify_Permissions(authID, set_conversation) {
		return "", Invalid_Permissions()
	}

	assistant_id, err := create_assistant(instructions, conversation_id)
	if err != nil {
		return "", Error_With_External_Service()
	}

	run, err := create_conversation(assistant_id, message)
	if err != nil {
		delete_assistant(assistant_id)
		return "", Error_With_External_Service()
	}

	var conversation Conversation

	conversation.ConversationID = conversation_id
	conversation.AssistantID = assistant_id
	conversation.ThreadID = run.ThreadID
	conversation.RunID = run.ID

	user, err := retrieve_user_auth_token(authID)
	if err != nil {
		delete_assistant(conversation.AssistantID)
		delete_thread(conversation.ThreadID)
		log.Println("Could not find user:", err)
		return "Could not find user", File_Error()
	}

	err = create_conversation_record(user.Username, conversation)
	if err != nil {
		delete_assistant(conversation.AssistantID)
		delete_thread(conversation.ThreadID)
		log.Println("Conversation not created:", err)
		return "Conversation not created", File_Error()
	}

	return get_last_message(conversation)
}

// TODO: finish Continue persona conversation function
func Continue_Persona_Conversation(authID, message, conversation_id string) (string, error) {
	if !Verify_Permissions(authID, set_conversation) {
		return "", Invalid_Permissions()
	}

	user, err := retrieve_user_auth_token(authID)
	if err != nil {
		log.Println("User not found:", err)
		return "User not found", Invalid_Data()
	}

	conversation, err := get_conversation(user.Username, conversation_id)
	if err != nil {
		log.Println("Message not found:", err)
		return "Message not found", File_Error()
	}

	run, err := update_conversation(conversation.AssistantID, conversation.ThreadID, message)
	if err != nil {
		log.Println("Thread could not be updated:", err)
		return "Thread could not be updated", Error_With_External_Service()
	}

	conversation.RunID = run.ID

	err = update_conversation_run_id(user.Username, conversation)
	if err != nil {
		log.Println("Converstation not updated:", err)
		return "Converstation not updated", File_Error()
	}

	return get_last_message(conversation)
}

// TODO: finish end persona conversation function
func End_Persona_Conversation(authID, conversation_id string) error {
	if !Verify_Permissions(authID, set_conversation) {
		return Invalid_Permissions()
	}

	user, err := retrieve_user_auth_token(authID)
	if err != nil {
		log.Println("User not found:", err)
		return Invalid_Data()
	}

	conversation, err := get_conversation(user.Username, conversation_id)
	if err != nil {
		log.Println("Conversation not found:", err)
		return File_Error()
	}

	err = delete_conversation(conversation)
	if err != nil {
		return Error_With_External_Service()
	}

	err = delete_conversation_record(user.Username, conversation_id)
	if err != nil {
		log.Println("Conversation object not deleted:", err)
		return File_Error()
	}

	return nil
}

// TODO: finish end all persona conversations function
func End_All_Persona_Conversations(authID string) error {
	if !Verify_Permissions(authID, set_conversation) {
		return Invalid_Permissions()
	}

	user, err := retrieve_user_auth_token(authID)
	if err != nil {
		log.Println("User not found:", err)
		return Invalid_Data()
	}

	conversations, err := get_all_conversations(user.Username)
	if err != nil {
		log.Println("Could not get conversation file:", err)
		return File_Error()
	}

	for _, conversation := range conversations {
		err = delete_conversation(conversation)
		if err != nil {
			log.Println("Error deleting Conversation:", err)
			return Error_With_External_Service()
		}
	}

	err = delete_conversation_file(user.Username)
	if err != nil {
		log.Println("Error deleting file:", err)
		return File_Error()
	}

	return nil
}
