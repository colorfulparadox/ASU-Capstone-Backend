package database

import (
	"context"
	"log"
	"math/rand/v2"
	"os"

	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	ConversationID string `json:"conversation_id"`
	AssistantID    string `json:"assistant_id"`
	ThreadID       string `json:"thread_id"`
	RunID          string `json:"run_id"`
}

func New_Persona(persona_name, ai_name, description, instructions string) error {
	persona := Persona{
		Name:         persona_name,
		AIName:       ai_name,
		Description:  description,
		Instructions: instructions,
	}

	return create_persona(persona)
}

func Conversation_Test(authID, conversation_id string) string {
	var conversation Conversation
	conversation.ConversationID = conversation_id
	conversation.AssistantID = "assistant id"
	conversation.ThreadID = "thread id"
	conversation.RunID = "run id"

	err := create_conversation(authID, conversation)

	return err.Error()
}

// TODO: Finish function and possibly add ability to select specific persona model
func Start_Persona_Conversation(authID, message, conversation_id string) (string, error) {
	if !Verify_Permissions(authID, start_conversation) {
		return "", Invalid_Permissions()
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Chooses a random persona to use
	personas, err := retrieve_persona_list()
	if err != nil {
		return "", err
	}
	persona_value := rand.IntN(len(personas))
	persona := personas[persona_value]

	var thread openai.ThreadRequest
	var run openai.CreateThreadAndRunRequest

	thread.Messages = []openai.ThreadMessage{
		{
			Role:    "user",
			Content: message,
		},
	}

	run.AssistantID = persona.AssistantID
	run.Thread = thread

	run_return, err := client.CreateThreadAndRun(context.Background(), run)

	if err != nil {
		log.Println("Thread not created:", err)
		return "Thread not created", err
	}

	var conversation Conversation

	conversation.ConversationID = conversation_id
	conversation.AssistantID = persona.AssistantID
	conversation.ThreadID = run_return.ThreadID
	conversation.RunID = run_return.ID

	err = create_conversation(authID, conversation)
	if err != nil {
		log.Println("Conversation not created:", err)
		return "Conversation not created", File_Error()
	}

	return get_last_message(conversation)
}

// TODO: finish Continue persona conversation function
func Continue_Persona_Conversation(authID, message, conversation_id string) (string, error) {
	if !Verify_Permissions(authID, continue_conversation) {
		return "", Invalid_Permissions()
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	conversation, err := get_conversation(authID, conversation_id)
	if err != nil {
		log.Println("Message not found:", err)
		return "Message not found", File_Error()
	}

	var message_request openai.MessageRequest
	var run_request openai.RunRequest

	message_request.Role = "user"
	message_request.Content = message

	_, err = client.CreateMessage(context.Background(), conversation.ThreadID, message_request)
	if err != nil {
		log.Println("Message not created:", err)
		return "Message not created", err
	}

	run_request.AssistantID = conversation.AssistantID

	run_return, err := client.CreateRun(context.Background(), conversation.ThreadID, run_request)
	if err != nil {
		log.Println("Message not created:", err)
		return "Message not created", err
	}

	conversation.RunID = run_return.ID

	err = update_conversation_run_id(authID, conversation)
	if err != nil {
		log.Println("Converstation not updated:", err)
		return "Converstation not updated", File_Error()
	}

	return get_last_message(conversation)
}

// TODO: finish end persona conversation function
func End_Persona_Conversation(authID, conversation_id string) error {
	if !Verify_Permissions(authID, end_conversation) {
		return Invalid_Permissions()
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	conversation, err := get_conversation(authID, conversation_id)
	if err != nil {
		log.Println("Conversation not found:", err)
		return File_Error()
	}

	_, err = client.DeleteThread(context.Background(), conversation.ThreadID)
	if err != nil {
		log.Println("Thread not deleted:", err)
		return Error_With_External_Service()
	}

	err = delete_conversation(authID, conversation_id)
	if err != nil {
		log.Println("Conversation object not deleted:", err)
		return File_Error()
	}

	return nil
}

// TODO: finish end all persona conversations function
func End_All_Persona_Conversations(authID string) error {
	if !Verify_Permissions(authID, end_conversation) {
		return Invalid_Permissions()
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	conversations, err := get_all_conversations(authID)
	if err != nil {
		log.Println("Could not get conversation file:", err)
		return File_Error()
	}

	for _, conversation := range conversations {
		_, err := client.DeleteThread(context.Background(), conversation.ThreadID)
		if err != nil {
			log.Println("Error deleting Thread:", err)
			return Error_With_External_Service()
		}
	}

	err = delete_conversation_file(authID)
	if err != nil {
		log.Println("Error deleting file:", err)
		return File_Error()
	}

	return nil
}
