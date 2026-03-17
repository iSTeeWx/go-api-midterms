package main

import (
	"examblanc/app"
	"examblanc/db"
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hey")
}

func main() {

	db.Instance = db.Connect()

	http.HandleFunc("GET /{$}", index)

	http.HandleFunc("GET /athletes/{$}", app.GetAthletes)
	http.HandleFunc("POST /athletes/{$}", app.PostAthletes)
	http.HandleFunc("GET /athletes/{id}/{$}", app.GetAthlete)
	http.HandleFunc("PUT /athletes/{id}/{$}", app.PutAthlete)
	http.HandleFunc("DELETE /athletes/{id}/{$}", app.DeleteAthlete)

	const host = "localhost:8080"

	fmt.Printf("Starting server on %v\n", host)
	http.ListenAndServe(host, nil)
}
