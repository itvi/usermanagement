package main

import (
	"database/sql"
	"log"
	"net/http"
	"umanagement/cmd/web/handler"
	"umanagement/pkg/model/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := openDB("./user.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	app := &handler.Application{
		Handlers: &handler.Handlers{
			User: &handler.UserHandler{M: &sqlite.UserModel{DB: db}},
		},
	}
	server := &http.Server{
		Addr:    ":9999",
		Handler: app.Routes(),
	}

	log.Println("Starting...")
	log.Fatal(server.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
