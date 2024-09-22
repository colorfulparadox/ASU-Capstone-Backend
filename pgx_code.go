package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func pgx_example_code() {
	databaseUrl := "postgres://project-persona:T%7D%3F_%5D0Lu8I98@postgres.blusnake.net:35432/project-persona"

	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	//CREATING TABLES EXAMPLE====================================================================================================================

	// Create the auth_tokens table
	_, err = conn.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS auth_tokens (
        id SERIAL PRIMARY KEY,
        date_issued TIMESTAMPTZ NOT NULL
    )`)
	if err != nil {
		log.Fatalf("Failed to create auth_tokens table: %v\n", err)
	}

	// Create the users table
	_, err = conn.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        emp_id TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        points INTEGER NOT NULL DEFAULT 0,
        permission_level INTEGER NOT NULL,
        email TEXT UNIQUE NOT NULL,
        auth_token_id INTEGER,
        date_expr TIMESTAMPTZ,
        FOREIGN KEY (auth_token_id) REFERENCES auth_tokens(id)
    )`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v\n", err)
	}

	log.Println("Tables created successfully")

	//CREATING DATA EXAMPLE====================================================================================================================

	// Transaction is useful if you're inserting into multiple tables
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Insert into auth_tokens table first
	var authTokenID int
	dateIssued := time.Now()

	err = tx.QueryRow(context.Background(),
		"INSERT INTO auth_tokens (date_issued) VALUES ($1) RETURNING id", dateIssued).Scan(&authTokenID)

	if err != nil {
		log.Fatalf("Failed to insert into auth_tokens: %v", err)
	}

	// Now insert into users table
	name := "John Doe"
	empID := "12345"
	username := "jdoe"
	password := "securepassword"
	points := 100
	permissionLevel := 1
	email := "jdoe@example.com"
	dateExpr := time.Now().Add(30 * 24 * time.Hour) // 30 days expiry

	_, err = tx.Exec(context.Background(),
		`INSERT INTO users 
        (name, emp_id, username, password, points, permission_level, email, auth_token_id, date_expr) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		name, empID, username, password, points, permissionLevel, email, authTokenID, dateExpr)

	if err != nil {
		log.Fatalf("Failed to insert into users: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Inserted user with auth token")

	//EDITING DATA EXAMPLE====================================================================================================================

	currentEmpID := "12345"
	newName := "Jane Doe"
	newPoints := 200

	// Perform the update
	commandTag, err := conn.Exec(context.Background(),
		"UPDATE users SET name = $1, points = $2 WHERE emp_id = $3",
		newName, newPoints, currentEmpID)

	if err != nil {
		log.Fatalf("Failed to update row: %v\n", err)
	}

	// Check how many rows were affected (should be 1 if successful)
	if commandTag.RowsAffected() != 1 {
		log.Fatalf("No row found with emp_id = %s", empID)
	}

	fmt.Println("User updated successfully")
}
