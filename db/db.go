package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	PORT = 5432
)

// ErrNoMatch is returned when we request a row that doesn't exist
var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.DB
}

func Initialize(username, password, database string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
	connected := false
	for !connected {
		conn, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Could not set up database: %v, retrying.", err)
			time.Sleep(3 * time.Second)
			continue
		}
		db.Conn = conn
		err = db.Conn.Ping()
		if err != nil {
			log.Printf("Could not set up database: %v, retrying.", err)
			time.Sleep(3 * time.Second)
			continue
		}
		connected = true
	}
	log.Println("Database connection established")

	driver, err := postgres.WithInstance(db.Conn, &postgres.Config{})
	if err != nil {
		return db, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://home/db/migrations",
		"noob_db", driver)
	if err != nil {
		fmt.Println("lol")
		return db, err
	}
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	log.Println("Migrations running")

	return db, nil
}
