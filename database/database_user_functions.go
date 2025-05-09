package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// User struct for holding users while working on them
type User struct {
	UserID           int       `json:"id"`
	Name             string    `json:"name"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	PasswordHash     string    `json:"password_hash"`
	Sentiment_Points int       `json:"sentiment_points"`
	Sales_Points     int       `json:"sales_points"`
	Knowledge_Points int       `json:"knowledge_points"`
	PermissionLevel  int       `json:"permission_level"`
	Email            string    `json:"email"`
	AuthToken        string    `json:"auth_token"`
	DateIssued       time.Time `json:"date_issued"`
	DateExpr         time.Time `json:"date_expr"`
}

// The table layout for users inside postgres
const usersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    sentiment_points INT NOT NULL DEFAULT 0,
    sales_points INT NOT NULL DEFAULT 0,
    knowledge_points INT NOT NULL DEFAULT 0,
    permission_level INT NOT NULL DEFAULT 0,
    email VARCHAR(255) UNIQUE NOT NULL,
    auth_token TEXT UNIQUE NOT NULL,
    date_issued TIMESTAMP NOT NULL,
    date_expr TIMESTAMP NOT NULL
);`

//Function for creating tables ====================================================================

func create_users_table() {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	// Create tables (if they don't exist)
	_, err = conn.Exec(context.Background(), usersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	user := User{
		Name:            "John Doe",
		Username:        "johndoe",
		Password:        "securepassword",
		PermissionLevel: 1,
		Email:           "john.doe@example.com",
	}

	create_user(user)
}

//Function for adding user data including auth token===============================================

// takes user object returns if user was created or not
func create_user(user User) error {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	username_user, email_user, _ := retrieve_user_pass_conn(conn, user.Username, user.Email)

	log.Println("Checking User")

	if username_user.Username == user.Username {
		err = errors.New("Username " + user.Username + " already in use")
		log.Println(err)
		return err
	}

	if email_user.Email == user.Email {
		err = errors.New("Email " + user.Email + " already in use")
		log.Println(err)
		return err
	}

	log.Printf("Creating user\n")

	user.AuthToken = GenerateUUID()
	var userID int
	// Prepare the SQL statement
	insertSQL := `INSERT INTO users (name, username, password, password_hash, permission_level, email, auth_token, date_issued, date_expr)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`

	// Execute the SQL statement using a prepared statement
	err = conn.QueryRow(context.Background(), insertSQL,
		user.Name, user.Username, user.Password, HashPassword(user.Password), user.PermissionLevel, user.Email, user.AuthToken, time.Now(), time.Now().AddDate(0, 0, 0)).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}

	randomize_auth_token(user.AuthToken)

	log.Printf("User created with ID: %d\n", userID)
	return nil
}

//Function for editing user data ==================================================================

func update_user(username string, user User) error {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	username_user, email_user, err := retrieve_user_pass_conn(conn, username, "")
	if err != nil {
		return err
	}

	// This checks to see if username is connected to a real user
	if username_user.Username != username {
		err = errors.New("User: " + user.Username + " does not exist")
		log.Println(err)
		return err
	}
	//This checks to make sure the desired info is not already in use
	username_user, email_user, err = retrieve_user_pass_conn(conn, user.Username, user.Email)
	if err != nil {
		return err
	}
	if username_user.UserID != user.UserID {
		log.Println(username_user.UserID, user.UserID)
		if username_user.Username == user.Username {
			err = errors.New("Username: " + user.Username + " already in use")
			log.Println(err)
			return err
		}
	}
	if email_user.UserID != user.UserID {
		if email_user.Email == user.Email {
			err = errors.New("Email: " + user.Username + " already in use")
			log.Println(err)
			return err
		}
	}

	if user.Password == "" {
		// Prepare the SQL statement for updating the user's name
		updateUserSQL := `UPDATE users SET name = $1, username = $2, sentiment_points = $3, sales_points = $4, knowledge_points = $5, permission_level = $6, email = $7
			WHERE username = $8;`

		// Execute the SQL statement using a prepared statement
		_, err = conn.Exec(context.Background(), updateUserSQL, user.Name, user.Username, user.Sentiment_Points, user.Sales_Points, user.Knowledge_Points, user.PermissionLevel, user.Email, username)
		if err != nil {
			log.Fatalf("Failed to update user's info: %v\n", err)
			return err
		}
	} else {
		// Prepare the SQL statement for updating the user's name
		updateUserSQL := `UPDATE users SET name = $1, username = $2, password = $3, password_hash = $4, sentiment_points = $5, sales_points = $6, knowledge_points = $7, permission_level = $8, email = $9
			WHERE username = $10;`

		// Execute the SQL statement using a prepared statement
		_, err = conn.Exec(context.Background(), updateUserSQL, user.Name, user.Username, user.Password, HashPassword(user.Password), user.Sentiment_Points, user.Sales_Points, user.Knowledge_Points, user.PermissionLevel, user.Email, username)
		if err != nil {
			log.Fatalf("Failed to update user's info: %v\n", err)
			return err
		}
	}
	log.Println("Update Complete")
	return nil
}

//Function for deleting a user ====================================================================

// takes user object returns if user was created or not
func delete_user(username string) error {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	username_user, _, err := retrieve_user_pass_conn(conn, username, "")
	if err != nil {
		return err
	}
	if username_user.Username == "" {
		err = errors.New("User: " + username + " does not exist")
		log.Println(err)
		return err
	}

	// Prepare the SQL statement
	deleteUserSQL := `DELETE FROM users 
	WHERE username = $1;`

	// Execute the SQL statement using a prepared statement
	_, err = conn.Exec(context.Background(), deleteUserSQL, username)
	if err != nil {
		log.Fatalf("Failed to delete user: %v\n", err)
		return err
	}

	log.Printf("Deleting User: %s\n", username)
	return nil
}

//Functions for retrieving user data ==============================================================

func retrieve_user_username(username string) (user User, err error) {
	conn := establish_connection()
	defer conn.Close()
	defer log.Println("Conn Closed")

	// Prepare the SQL statement for selecting the user's data
	selectUserSQL := `SELECT id, name, username, password, password_hash, sentiment_points, sales_points, knowledge_points, permission_level, email, auth_token, date_issued, date_expr
    	FROM users
    	WHERE username = $1;`

	err = conn.QueryRow(context.Background(), selectUserSQL, username).Scan(
		&user.UserID, //the id variable should not be used outside the backend
		&user.Name,
		&user.Username,
		&user.Password,
		&user.PasswordHash,
		&user.Sentiment_Points,
		&user.Sales_Points,
		&user.Knowledge_Points,
		&user.PermissionLevel,
		&user.Email,
		&user.AuthToken,
		&user.DateIssued,
		&user.DateExpr,
	)
	if err != nil {
		log.Printf("Failed to retrieve user: %s, %v\n", username, err)
	} else {
		log.Printf("Retrieved user: %s\n", user.Username)
	}
	return
}

func retrieve_user_auth_token(auth_token string) (user User, err error) {
	conn := establish_connection()
	defer conn.Close()
	defer log.Println("Conn Closed")

	// Prepare the SQL statement for selecting the user's data
	selectUserSQL := `SELECT id, name, username, password, sentiment_points, sales_points, knowledge_points, permission_level, email, auth_token, date_issued, date_expr
    	FROM users
    	WHERE auth_token = $1;`

	err = conn.QueryRow(context.Background(), selectUserSQL, auth_token).Scan(
		&user.UserID, //the id variable should not be used outside the backend
		&user.Name,
		&user.Username,
		&user.Password,
		&user.Sentiment_Points,
		&user.Sales_Points,
		&user.Knowledge_Points,
		&user.PermissionLevel,
		&user.Email,
		&user.AuthToken,
		&user.DateIssued,
		&user.DateExpr,
	)
	if err != nil {
		log.Printf("Failed to retrieve user: %s, %v\n", user.Username, err)
	} else {
		log.Printf("Retrieved user: %s\n", user.Username)
	}
	return
}

//private version of the Rerieve_user function that uses conn and err so a new connection does not have to be made

func retrieve_user_pass_conn(conn *pgxpool.Pool, username, email string) (username_user, email_user User, err error) {
	// Prepare the SQL statement for selecting the user's data

	selectUserSQL := `SELECT id, name, username, password, sentiment_points, sales_points, knowledge_points, permission_level, email, auth_token, date_issued, date_expr
    	FROM users
    	WHERE username = $1;`

	err = conn.QueryRow(context.Background(), selectUserSQL, username).Scan(
		&username_user.UserID, //the id variable should not be used outside the backend
		&username_user.Name,
		&username_user.Username,
		&username_user.Password,
		&username_user.Sentiment_Points,
		&username_user.Sales_Points,
		&username_user.Knowledge_Points,
		&username_user.PermissionLevel,
		&username_user.Email,
		&username_user.AuthToken,
		&username_user.DateIssued,
		&username_user.DateExpr,
	)
	if err != nil {
		log.Printf("No users with that name exist: %v\n", err)
	}

	if email != "" {

		selectUserSQL = `SELECT id, name, username, password, sentiment_points, sales_points, knowledge_points, permission_level, email, auth_token, date_issued, date_expr
    		FROM users
    		WHERE email = $1;`

		err = conn.QueryRow(context.Background(), selectUserSQL, email).Scan(
			&email_user.UserID, //the id variable should not be used outside the backend
			&email_user.Name,
			&email_user.Username,
			&email_user.Password,
			&email_user.Sentiment_Points,
			&email_user.Sales_Points,
			&email_user.Knowledge_Points,
			&email_user.PermissionLevel,
			&email_user.Email,
			&email_user.AuthToken,
			&email_user.DateIssued,
			&email_user.DateExpr,
		)
		if err != nil {
			log.Printf("No users with that email exist: %v\n", err)
		}
	}
	log.Println(username_user.UserID, email_user.UserID)
	return
}

//Function for retrieving a list of all users =====================================================

func retrieve_user_list() ([]User, error) {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	var userList []User
	var current_user int

	getIDsSQL := `SELECT id FROM users`
	validIDs, err := conn.Query(context.Background(), getIDsSQL)
	if err != nil {
		log.Printf("Could not get user ID's\n")
		log.Printf("Returned error was: %v\n", err)
		return nil, err
	}

	for validIDs.Next() {
		var user User
		validIDs.Scan(&current_user)
		// Prepare the SQL statement for selecting the user's data
		selectUserSQL := `SELECT id, name, username, password, sentiment_points, sales_points, knowledge_points, permission_level, email, auth_token, date_issued, date_expr
			FROM users
			WHERE id = $1;`

		err = conn.QueryRow(context.Background(), selectUserSQL, current_user).Scan(
			&user.UserID, //the id variable should not be used outside the backend
			&user.Name,
			&user.Username,
			&user.Password,
			&user.Sentiment_Points,
			&user.Sales_Points,
			&user.Knowledge_Points,
			&user.PermissionLevel,
			&user.Email,
			&user.AuthToken,
			&user.DateIssued,
			&user.DateExpr,
		)
		if err != nil {
			log.Printf("An Error occured retrieving user data: %v\n", err)
			break
		} else {
			log.Printf("Retrieved user: %s\n", user.Username)
			userList = append(userList, user)
			continue
		}
	}

	return userList, nil
}

//Function for randomizing auth token ============================================================

func randomize_auth_token(auth_token string) error {
	conn := establish_connection()
	var err error
	defer conn.Close()
	defer log.Println("Conn Closed")

	updateNameSQL := `UPDATE users SET auth_token = $1, date_issued = $2, date_expr = $3
		WHERE auth_token = $4;`

	dateIssued := time.Now().UTC()

	local_hour := time.Now().Hour()
	hour_offset := 23 - local_hour

	expires := time.Now().Add(time.Hour*time.Duration(hour_offset)).AddDate(0, 0, Authentication_Token_Forced_Time_Reset).UTC()

	token := GenerateUUID()

	log.Println("Randomizing Authentication Token")

	// Execute the SQL statement using a prepared statement
	_, err = conn.Exec(context.Background(), updateNameSQL, token, dateIssued, expires, auth_token)
	if err != nil {
		log.Fatalf("Failed to randomize user's auth_token: %v\n", err)
		return err
	}

	return nil
}
