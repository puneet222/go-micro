package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"
const maxTries = 10
const backoffTime = 2

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// TODO connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Database connection failed!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDBConnection(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	tries := 0
	for {
		conn, err := openDBConnection(dsn)
		if err == nil {
			log.Printf("Connected to Postgres database!")
			return conn
		}
		log.Printf("Postgres database not ready yet: %v\n", err)
		tries++

		if tries > maxTries {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			return nil
		}

		time.Sleep(time.Second * backoffTime)
		fmt.Printf("Backing off for %d seconds", backoffTime)
	}
}
