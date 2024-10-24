package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

//Function for setting up connection ==============================================================

func establish_connection() (conn *pgxpool.Pool) {
	// Set up connection to the PostgreSQL server
	var err error

	var dsn string

	host := "/cloudsql/" + os.Getenv("CLOUD_SQL_CONNECTION_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	if os.Getenv("MODE") == "release" {
		dsn = "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	} else {
		host := os.Getenv("CLOUD_SQL_CONNECTION_NAME")
		port := os.Getenv("DB_PORT")
		dsn = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	}
	conn, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("INTERNAL: Unable to connect to database: %v\n", err)
	}

	log.Println("Conn Opened")

	return
}
