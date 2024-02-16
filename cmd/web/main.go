package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// command-line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// router
	mux := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// routes
	mux.HandleFunc("/", dashboard)
	mux.HandleFunc("/projects", projectListView)
	mux.HandleFunc("/project/view", projectView)
	mux.HandleFunc("/project/create", projectCreate)

	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
