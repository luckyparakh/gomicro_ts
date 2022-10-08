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

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Print("Starting Auth Service")
	conn, err := connectToDB()
	if err != nil {
		log.Panic("Cannot connect to Postgres")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	for {
		db, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet...")
			counts++
		} else {
			log.Println("Postgres is ready.")
			return db, nil
		}

		if counts > 10 {
			log.Println("Postgres is not ready yet...leaving")
			return nil, err
		}
		log.Println("Backing off for 2 seconds")
		time.Sleep(2 * time.Second)
	}
}
