package database

import (
	"context"
	"errors"
	"log"
	"os"
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
func create_persona(new_persona Persona) bool {
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
			log.Println("Persona already exists:", err)
			return false
		}
	}

	ai, err := retrieve_ai(new_persona.AIName)
	if err != nil {
		log.Println("AI model could not be retrieved:", err)
		return false
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
		return false
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
	}

	log.Println(new_persona)

	return true
}

func retrieve_persona_pass_conn(conn *pgxpool.Pool, name string) (persona Persona, err error) {
	// Prepare the SQL statement for selecting the persona's data
	selectUserSQL := `SELECT id, name, ai_name, assistant_id, description, instructions, last_edit
			FROM personas
			WHERE name = $1;`

	err = conn.QueryRow(context.Background(), selectUserSQL, name).Scan(
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
		return
	}
	return
}

func retrieve_persona_list() []Persona {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	getIDsSQL := `SELECT id FROM users`
	validIDs, err := conn.Query(context.Background(), getIDsSQL)
	if err != nil {
		log.Printf("Could not get user ID's\n")
		log.Printf("Returned error was: %v\n", err)
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

	return personaList
}
