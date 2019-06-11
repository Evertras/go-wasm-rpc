package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Just serve everything in front/ to make it easy
	mux.Handle("/", http.FileServer(http.Dir("front")))

	log.Println("Starting on http://localhost:8000")

	s := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
