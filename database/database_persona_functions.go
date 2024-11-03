package database

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sashabaranov/go-openai"
)

type Persona struct {
	PersonaID    int       `json:"id"`
	Name         string    `json:"name"`
	AIName       string    `json:"ai_name"`
	AssistantID  string    `json:"assistant_id"`
	Description  string    `json:"description"`
	Instructions string    `json:"instructions"`
	LastEdit     time.Time `json:"last_edit"`
}

const personaTableSQL = `
CREATE TABLE IF NOT EXISTS personas (
    id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	ai_name VARCHAR(255) NOT NULL,
    assistant_id VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    instructions VARCHAR(255) NOT NULL,
	last_edit TIMESTAMP NOT NULL
);`

func create_persona_table() {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	// Create tables (if they don't exist)
	_, err = conn.Exec(context.Background(), personaTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	persona := Persona{
		Name:         "default",
		AIName:       "default",
		Description:  "A default AI created for testing",
		Instructions: "You are a default ai used for testing, your primary directive is to let the user know that you are functioning correctly",
	}

	create_persona(persona)
}

// TODO: finish create persona function
func create_persona(new_persona Persona) error {
	conn := establish_connection()
	defer conn.Close()
	defer log.Println("Conn Closed")
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	log.Println("Creating Persona")

	current_persona, err := retrieve_persona_pass_conn(conn, new_persona.Name)
	if err != nil {
		log.Println("AI model could not be retrieved:", err)
	} else {
		if current_persona.Name == new_persona.Name {
			err := errors.New("Persona already esists")
			return err
		}
	}

	ai, err := retrieve_ai(new_persona.AIName)
	if err != nil {
		log.Println("AI model could not be retrieved:", err)
		return err
	}

	var persona openai.AssistantRequest
	var file_search openai.AssistantTool
	// var vector_search openai.AssistantToolResource

	file_search.Type = openai.AssistantToolTypeFileSearch
	file_search.Function = nil

	// file_id := openai.AssistantToolFileSearch{
	// 	VectorStoreIDs: []string{ai.VectorID},
	// }

	// vector_search.FileSearch = &file_id

	persona.Model = Model_Name
	persona.Name = &new_persona.Name
	persona.Description = &new_persona.Description
	persona.Instructions = &new_persona.Instructions
	persona.Tools = append(persona.Tools, file_search)
	// persona.ToolResources = &vector_search

	persona_result, err := client.CreateAssistant(context.Background(), persona)
	if err != nil {
		log.Fatalln("Persona not stored in openai: ", err)
		return err
	}

	new_persona.AssistantID = persona_result.ID
	new_persona.LastEdit = time.Now()

	var personaID int

	insertSQL := `INSERT INTO personas (name, ai_name, assistant_id, description, instructions, last_edit)
    	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	// Execute the SQL statement using a prepared statement
	err = conn.QueryRow(context.Background(), insertSQL,
		new_persona.Name, new_persona.AIName, new_persona.AssistantID, new_persona.Description, new_persona.Instructions, ai.LastEdit).Scan(&personaID)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
		return err
	}

	log.Println(new_persona)

	return nil
}

func retrieve_persona_pass_conn(conn *pgxpool.Pool, name string) (Persona, error) {
	// Prepare the SQL statement for selecting the persona's data
	selectUserSQL := `SELECT id, name, ai_name, assistant_id, description, instructions, last_edit
			FROM personas
			WHERE name = $1;`

	var persona Persona

	err := conn.QueryRow(context.Background(), selectUserSQL, name).Scan(
		&persona.PersonaID, //the id variable should not be used outside the backend
		&persona.Name,
		&persona.AIName,
		&persona.AssistantID,
		&persona.Description,
		&persona.Instructions,
		&persona.LastEdit,
	)
	if err != nil {
		err = errors.New("Failed to retireve AI:" + err.Error())
		return persona, err
	}
	return persona, nil
}

func retrieve_persona_list() ([]Persona, error) {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	getIDsSQL := `SELECT id FROM users`
	validIDs, err := conn.Query(context.Background(), getIDsSQL)
	if err != nil {
		log.Printf("Could not get user ID's\n")
		log.Printf("Returned error was: %v\n", err)
		return nil, err
	}

	var personaList []Persona
	var current_persona int

	for validIDs.Next() {
		var persona Persona
		validIDs.Scan(&current_persona)
		// Prepare the SQL statement for selecting the user's data
		selectUserSQL := `SELECT id, name, ai_name, assistant_id, description, instructions, last_edit
			FROM personas
			WHERE id = $1;`

		err = conn.QueryRow(context.Background(), selectUserSQL, current_persona).Scan(
			&persona.PersonaID, //the id variable should not be used outside the backend
			&persona.Name,
			&persona.AIName,
			&persona.AssistantID,
			&persona.Description,
			&persona.Instructions,
			&persona.LastEdit,
		)
		if err != nil {
			log.Printf("An Error occured retrieving user data: %v\n", err)
			break
		} else {
			log.Printf("Retrieved user: %s\n", persona.Name)
			personaList = append(personaList, persona)
			continue
		}
	}

	return personaList, nil
}

func get_last_message(conversation Conversation) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var run_return openai.Run

	for i := 0; i < Completed_Timeout; i++ {
		time.Sleep(500 * time.Millisecond)
		log.Println("Time")
		run_return, err := client.RetrieveRun(context.Background(), conversation.ThreadID, conversation.RunID)
		if err != nil {
			log.Println("run not found: ", err)
			return "run not found", err
		}

		if run_return.Status == "completed" {
			break
		}
	}

	if run_return.Status != "completed" {
		log.Println("Timeout")
		return "Timeout", errors.New("Timeout")
	}

	limit := 1
	messages_return, err := client.ListMessage(context.Background(), conversation.ThreadID, &limit, nil, nil, nil, nil)
	if err != nil {
		log.Println("Messages not found: ", err)
		return "Messages not found", err
	}

	return messages_return.Messages[0].Content[0].Text.Value, nil
}

func get_conversation(authID, conversation_id string) (Conversation, error) {
	file, err := os.OpenFile(filepath.Join(Conversation_Path, authID+".json"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Temp file not found: ", err)
		return Conversation{}, err
	}
	defer file.Close()

	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		log.Println("Error decoding JSON (EOF error is expected when creating a new file):", err)
	}

	var conversation Conversation
	for _, current_conversation := range conversations {
		if current_conversation.ConversationID == conversation_id {
			conversation = current_conversation
		}
	}

	return conversation, nil
}

func create_conversation(authID string, conversation Conversation) error {
	conversations, err := get_all_conversations(authID)
	if err != nil {
		log.Println("Conversations not found: ", err)
		return err
	}

	file, err := os.OpenFile(filepath.Join(Conversation_Path, authID+".json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Temp file not created/found: ", err)
		return err
	}
	defer file.Close()

	for _, current_conversation := range conversations {
		if current_conversation.ConversationID == conversation.ConversationID {
			log.Println("Conversation already exists")
			return errors.New("conversation already exists")
		}
	}

	conversations = append(conversations, conversation)

	encoder := json.NewEncoder(file)
	encoder.Encode(conversations)

	return nil
}

// This is only to be used to change the RunID to work with openai systems DO NOT CHANGE ANYTHING ELSE WITH THIS FUNCTION
func update_conversation_run_id(authID string, conversation Conversation) error {
	conversations, err := get_all_conversations(authID)
	if err != nil {
		log.Println("Could not get conversation file:", err)
		return err
	}

	var position int
	var old_conversation Conversation
	for current_position, current_conversation := range conversations {
		if current_conversation.ConversationID == conversation.ConversationID {
			position = current_position
			old_conversation = current_conversation
		}
	}

	if old_conversation.ConversationID == "" {
		log.Println("Conversation not found")
		return Invalid_Data()
	}

	conversations[position] = conversation

	file, err := os.OpenFile(filepath.Join(Conversation_Path, authID+".json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("Temp file not created/found: ", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(conversations)

	return nil
}

func delete_conversation(authID string, conversation_id string) error {
	conversations, err := get_all_conversations(authID)
	if err != nil {
		log.Println("Could not get conversation file:", err)
		return err
	}

	var position int
	var conversation Conversation
	for current_position, current_conversation := range conversations {
		if current_conversation.ConversationID == conversation_id {
			position = current_position
			conversation = current_conversation
		}
	}

	if conversation.ConversationID == "" {
		log.Println("Conversation not found")
		return Invalid_Data()
	}

	if len(conversations) <= 1 {
		err = os.Remove(filepath.Join(Conversation_Path, authID+".json"))
		if err != nil {
			log.Println("Error deleting file:", err)
			return File_Error()
		}
	} else {
		var new_conversations []Conversation
		new_conversations = append(new_conversations, conversations[:position]...)
		new_conversations = append(new_conversations, conversations[position+1:]...)

		file, err := os.OpenFile(filepath.Join(Conversation_Path, authID+".json"), os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Println("Temp file not created/found: ", err)
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.Encode(new_conversations)
	}

	return nil
}

func get_all_conversations(authID string) ([]Conversation, error) {
	file, err := os.OpenFile(filepath.Join(Conversation_Path, authID+".json"), os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println("Temp file not found: ", err)
		return nil, err
	}
	defer file.Close()

	var conversations []Conversation
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conversations)
	if err != nil {
		log.Println("Error decoding JSON (EOF error is expected when creating a new file):", err)
	}

	return conversations, nil
}

func delete_conversation_file(authID string) error {
	file, err := os.OpenFile(filepath.Join(Conversation_Path, authID+".json"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Temp file not found: ", err)
		return File_Error()
	}
	defer file.Close()

	return nil
}
