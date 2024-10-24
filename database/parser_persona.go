package database

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	ConversationID string `json:"conversation_id"`
	AssistantID    string `json:"assistant_id"`
	ThreadID       string `json:"thread_id"`
	RunID          string `json:"run_id"`
}

func New_Persona(persona_name, ai_name, description, instructions string) int {
	persona := Persona{
		Name:         persona_name,
		AIName:       ai_name,
		Description:  description,
		Instructions: instructions,
	}

	if create_persona(persona) {
		return Success
	} else {
		return Invalid_Data
	}
}

func Get_Last_Message(conversation Conversation) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var run_return openai.Run

	for i := 0; i < Completed_Timeout; i++ {
		time.Sleep(500 * time.Millisecond)
		log.Println("Time")
		run_return, err := client.RetrieveRun(context.Background(), conversation.ThreadID, conversation.RunID)
		if err != nil {
			log.Println("run not found: ", err)
			return "run not found: " + err.Error()
		}

		if run_return.Status == "completed" {
			break
		}
	}

	if run_return.Status != "completed" {
		log.Println("Timeout")
		return "Timeout"
	}

	limit := 1
	messages_return, err := client.ListMessage(context.Background(), conversation.ThreadID, &limit, nil, nil, nil, nil)
	if err != nil {
		log.Println("Messages not found: ", err)
		return "Messages not found: " + err.Error()
	}

	return messages_return.Messages[0].Content[0].Text.Value
}

// TODO: Finish function and possibly add ability to select specific persona model
func Start_Persona_Conversation(authID, message, conversation_id string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	file, err := os.OpenFile(filepath.Join(Temp_Path, Threads_Path, authID+".json"), os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println("Temp file not created/found: ", err)
		return "Temp file not created: " + err.Error()
	}
	defer file.Close()

	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		log.Println("Error decoding JSON (EOF error is expected when creating a new file):", err)
		//return "Error decoding JSON:" + err.Error()
	}

	for _, conversation := range conversations {
		if conversation.ConversationID == conversation_id {
			log.Println("Error: Conversation already exists")
			return "Error: Conversation already exists"
		}
	}

	// Chooses a random persona to use
	personas := retrieve_persona_list()
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
		return "Thread not created: " + err.Error()
	}

	var conversation Conversation

	conversation.ConversationID = conversation_id
	conversation.AssistantID = persona.AssistantID
	conversation.ThreadID = run_return.ThreadID
	conversation.RunID = run_return.ID

	conversations = append(conversations, conversation)

	file, err = os.OpenFile(filepath.Join(Temp_Path, Threads_Path, authID+".json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Temp file not created/found: ", err)
		return "Temp file not created: " + err.Error()
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(conversations)

	return Get_Last_Message(conversation)
}

// TODO: finish Continue persona conversation function
func Continue_Persona_Conversation(authID, message, conversation_id string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	file, err := os.OpenFile(filepath.Join(Temp_Path, Threads_Path, authID+".json"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Temp file not found: ", err)
		return "Temp file not found: " + err.Error()
	}
	defer file.Close()

	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		log.Println("Error decoding JSON (EOF error is expected when creating a new file):", err)
		//return "Error decoding JSON:" + err.Error()
	}

	var conversation Conversation
	var conversation_value int

	for i, current_conversation := range conversations {
		if current_conversation.ConversationID == conversation_id {
			conversation = current_conversation
			conversation_value = i
			log.Println("Conversation found")
		}
	}

	if conversation.ConversationID == "" {
		log.Println("Conversation not found")
		return "Conversation not found"
	}

	var message_request openai.MessageRequest
	var run_request openai.RunRequest

	message_request.Role = "user"
	message_request.Content = message

	_, err = client.CreateMessage(context.Background(), conversation.ThreadID, message_request)
	if err != nil {
		log.Println("Message not created:", err)
		return "Message not created: " + err.Error()
	}

	run_request.AssistantID = conversation.AssistantID

	run_return, err := client.CreateRun(context.Background(), conversation.ThreadID, run_request)
	if err != nil {
		log.Println("Message not created:", err)
		return "Message not created: " + err.Error()
	}

	conversation.RunID = run_return.ID
	conversations[conversation_value].RunID = conversation.RunID

	file, err = os.OpenFile(filepath.Join(Temp_Path, Threads_Path, authID+".json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Temp file not created/found: ", err)
		return "Temp file not created: " + err.Error()
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(conversations)

	return Get_Last_Message(conversation)
}

// TODO: finish end persona conversation function
func End_Persona_Conversation(authID, conversation_id string) int {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	file, err := os.OpenFile(filepath.Join(Temp_Path, Threads_Path, authID+".json"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Temp file not found: ", err)
		return File_Error
	}
	defer file.Close()

	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		log.Println("Error decoding JSON (EOF error is expected when creating a new file):", err)
		//return "Error decoding JSON:" + err.Error()
	}

	var conversation Conversation

	for _, current_conversation := range conversations {
		if current_conversation.ConversationID == conversation_id {
			conversation = current_conversation
			log.Println("Conversation found")
		}
	}

	if conversation.ConversationID == "" {
		log.Println("Conversation not found")
		return Invalid_Data
	}

	_, err = client.DeleteThread(context.Background(), conversation.ThreadID)
	if err != nil {
		log.Println("Message not created:", err)
		return Error_With_External_Service
	}

	return Success
}

// TODO: finish end all persona conversations function
func End_All_Persona_Conversations(authID string) int {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	file, err := os.OpenFile(filepath.Join(Temp_Path, Threads_Path, authID+".json"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Temp file not found: ", err)
		return File_Error
	}
	defer file.Close()

	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		log.Println("Error decoding JSON (EOF error is expected when creating a new file):", err)
		return Invalid_Data
	}

	for _, conversation := range conversations {
		_, err := client.DeleteThread(context.Background(), conversation.ThreadID)
		if err != nil {
			log.Println("Error deleting Thread:", err)
			return Error_With_External_Service
		}
	}

	err = os.Remove(filepath.Join(Temp_Path, Threads_Path, authID+".json"))
	if err != nil {
		log.Println("Error deleting file:", err)
		return File_Error
	}

	return Success
}
