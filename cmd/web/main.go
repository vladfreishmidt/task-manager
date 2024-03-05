package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vladfreishmidt/task-manager/internal/models"
)

// application struct to hold the application-wide dependencies
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	projects       *models.ProjectModel
	workspaces     *models.WorkspaceModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// initialize a form decoder
	formDecoder := form.NewDecoder()

	// initialize a session manager
	sessionsManager := scs.New()
	sessionsManager.Store = mysqlstore.New(db)
	sessionsManager.Lifetime = 12 * time.Hour

	// application struct instance containing dependencies
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		projects:       &models.ProjectModel{DB: db},
		workspaces:     &models.WorkspaceModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionsManager,
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
