package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vladfreishmidt/task-manager/internal/models"
)

// application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	projects *models.ProjectModel
}

func main() {
	// command-line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:AF12redux@/db_task_manager?parseTime=true", "MySQL data source name")
	flag.Parse()

	// loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create db connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// application struct instance containing dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		projects: &models.ProjectModel{DB: db},
	}

	// new http.Server struct with configuration settings for the server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
