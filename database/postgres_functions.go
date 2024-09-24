package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User struct for holding users while working on them
type User struct {
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	PasswordHash    string    `json:"password_hash"`
	Points          int       `json:"points"`
	PermissionLevel int       `json:"permission_level"`
	Email           string    `json:"email"`
	AuthToken       string    `json:"auth_token"`
	DateIssued      time.Time `json:"date_issued"`
	DateExpr        time.Time `json:"date_expr"`
}

// The table layout for users inside postgres
const usersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    points INT NOT NULL DEFAULT 0,
    permission_level INT NOT NULL DEFAULT 0,
    email VARCHAR(255) UNIQUE NOT NULL,
    auth_token TEXT UNIQUE NOT NULL,
    date_issued TIMESTAMP NOT NULL,
    date_expr TIMESTAMP NOT NULL
);`

// Database url
const databaseUrl = "postgres://project-persona:T%7D%3F_%5D0Lu8I98@postgres.blusnake.net:35432/project-persona"

// func pgx_examples() {
// 	create_users_table()

// 	user := User{
// 		Name:            "John Doe",
// 		Username:        "johndoe",
// 		Password:        "securepassword",
// 		Points:          rand.IntN(1000),
// 		PermissionLevel: 0,
// 		Email:           "john.doe@example.com",
// 	}

// 	create_user(user)

// 	username := "johndoe"
// 	user = User{
// 		Name:            "Jane Doe",
// 		Username:        "janedoe",
// 		Password:        "securepassword",
// 		Points:          100,
// 		PermissionLevel: 1,
// 		Email:           "jane.doe@example.com",
// 	}

// 	update_user(username, user)

// 	retrieve_user_username("Username")

// 	retrieve_user_auth_token("Token")

// 	randomize_auth_token("Token")
// }

//Function for setting up connection ==============================================================

func establish_connection() (conn *pgxpool.Pool) {
	// Set up connection to the PostgreSQL server
	var err error
	conn, err = pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("INTERNAL: Unable to connect to database: %v\n", err)
	}
	return
}

//Function for creating tables ====================================================================

func create_users_table() {
	conn := establish_connection()
	var err error
	defer conn.Close()

	// Create tables (if they don't exist)
	_, err = conn.Exec(context.Background(), usersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	user := User{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "securepassword",
		Points:          0,
		PermissionLevel: 1,
		Email:           "john.doe@example.com",
	}

	create_user(user)
}

//Function for adding user data including auth token===============================================

// takes user object returns if user was created or not
func create_user(user User) bool {
	conn := establish_connection()
	var err error
	defer conn.Close()

	compare_user := retrieve_user_username_pass_conn(conn, user.Username)
	if compare_user.Username == user.Username {
		log.Printf("Username: %s already in use\n", user.Username)
		return false
	} else {
		log.Printf("Creating user\n")
	}

	user.AuthToken = GenerateUUID()

	// Prepare the SQL statement
	insertSQL := `INSERT INTO users (name, username, password, password_hash, points, permission_level, email, auth_token, date_issued, date_expr)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;`

	// Execute the SQL statement using a prepared statement
	var userID int
	err = conn.QueryRow(context.Background(), insertSQL,
		user.Name, user.Username, user.Password, HashPassword(user.Password), user.Points, user.PermissionLevel, user.Email, user.AuthToken, time.Now(), time.Now().AddDate(0, 0, 0)).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

	randomize_auth_token(user.AuthToken)

	log.Printf("User inserted with ID: %d\n", userID)

	return true
}

//Function for editing user data ==================================================================

func update_user(username string, user User) bool {
	conn := establish_connection()
	var err error
	defer conn.Close()

	// This checks to see if username is connected to a real user
	if retrieve_user_username_pass_conn(conn, username).Username != username {
		log.Printf("User: %s does not exist\n", username)
		return false
	}
	//This checks to make sure the desired info is not already in use
	compare_user := retrieve_user_username_pass_conn(conn, user.Username)
	if compare_user.Username == user.Username {
		log.Printf("Username: %s already in use\n", user.Username)
		return false
	}

	if compare_user.Email == user.Email {
		log.Printf("Email: %s already in use\n", user.Email)
		return false
	}

	// Prepare the SQL statement for updating the user's name
	updateNameSQL := `UPDATE users SET name = $1, username = $2, password = $3, password_hash = $4, points = $5, permission_level = $6, email = $7
	WHERE username = $8;`

	// Execute the SQL statement using a prepared statement
	_, err = conn.Exec(context.Background(), updateNameSQL, user.Name, user.Username, user.Password, HashPassword(user.Password), user.Points, user.PermissionLevel, user.Email, username)
	if err != nil {
		log.Fatalf("Failed to update user's info: %v\n", err)
		return false
	}

	return true

	//log.Printf("User: %s, New Name: %s New Username: %s\n", username, user.Name, user.Username)
}

//Functions for retrieving user data ==============================================================

func retrieve_user_username(username string) (user User) {
	conn := establish_connection()
	var err error
	defer conn.Close()

	// Prepare the SQL statement for selecting the user's data
	// the id column is just here for completeness and should not be referenced in actual deployment
	selectUserSQL := `SELECT id, name, username, password, password_hash, points, permission_level, email, auth_token, date_issued, date_expr
    FROM users
    WHERE username = $1;`

	var current_userID int
	err = conn.QueryRow(context.Background(), selectUserSQL, username).Scan(
		&current_userID, //the id variable is here for completeness and should not be referenced in actual deployment
		&user.Name,
		&user.Username,
		&user.Password,
		&user.PasswordHash,
		&user.Points,
		&user.PermissionLevel,
		&user.Email,
		&user.AuthToken,
		&user.DateIssued,
		&user.DateExpr,
	)
	if err != nil {
		log.Printf("Failed to retrieve user data from username: %v\n", err)
	}

	//log.Printf("User Data: %+v\n", user)

	return
}

func retrieve_user_auth_token(auth_token string) (user User) {
	conn := establish_connection()
	var err error
	defer conn.Close()

	// Prepare the SQL statement for selecting the user's data
	// the id column is just here for completeness and should not be referenced in actual deployment
	selectUserSQL := `SELECT id, name, username, password, points, permission_level, email, auth_token, date_issued, date_expr
    FROM users
    WHERE auth_token = $1;`

	var current_userID int
	err = conn.QueryRow(context.Background(), selectUserSQL, auth_token).Scan(
		&current_userID, //the id variable is here for completeness and should not be referenced in actual deployment
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Points,
		&user.PermissionLevel,
		&user.Email,
		&user.AuthToken,
		&user.DateIssued,
		&user.DateExpr,
	)
	if err != nil {
		log.Printf("Failed to retrieve user data from auth_token: %v\n", err)
	}

	//log.Printf("User Data: %+v\n", user)

	return
}

//private version of the Rerieve_user function that uses conn and err so a new connection does not have to be made

func retrieve_user_username_pass_conn(conn *pgxpool.Pool, username string) (user User) {
	// Prepare the SQL statement for selecting the user's data
	// the id column is just here for completeness and should not be referenced in actual deployment
	selectUserSQL := `SELECT id, name, username, password, points, permission_level, email, auth_token, date_issued, date_expr
    FROM users
    WHERE username = $1;`

	var current_userID int
	err := conn.QueryRow(context.Background(), selectUserSQL, username).Scan(
		&current_userID, //the id variable is here for completeness and should not be referenced in actual deployment
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Points,
		&user.PermissionLevel,
		&user.Email,
		&user.AuthToken,
		&user.DateIssued,
		&user.DateExpr,
	)
	if err != nil {
		log.Printf("No users with that name exist: %v\n", err)
	}

	//log.Printf("User Data: %+v\n", user)

	return
}

//Function for randomizing auth token ============================================================

func randomize_auth_token(auth_token string) {
	conn := establish_connection()
	var err error
	defer conn.Close()

	updateNameSQL := `UPDATE users SET auth_token = $1, date_issued = $2, date_expr = $3
	WHERE auth_token = $4;`

	dateIssued := time.Now().UTC()
	expires := time.Now().AddDate(0, 0, 7).UTC()

	token := GenerateUUID()

	// Execute the SQL statement using a prepared statement
	_, err = conn.Exec(context.Background(), updateNameSQL, token, dateIssued, expires, auth_token)
	if err != nil {
		log.Fatalf("Failed to randomize user's auth_token: %v\n", err)
	}
}
