package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// command-line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// application struct instance containing dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// new http.Server struct with configuration settings for the server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
