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

	http.HandleFunc("GET /judges/{$}", app.GetJudges)
	http.HandleFunc("POST /judges/{$}", app.PostJudge)
	http.HandleFunc("GET /judges/{id}/{$}", app.GetJudge)
	http.HandleFunc("DELETE /judges/{id}/{$}", app.DeleteJudge)

	http.HandleFunc("POST /matches/{$}", app.PostMatch)
	http.HandleFunc("GET /judges/{id}/matches/{$}", app.GetMatchesOfJudge)
	http.HandleFunc("DELETE /matches/{id}/{$}", app.DeleteMatch)

	const host = "localhost:8080"

	fmt.Printf("Starting server on %v\n", host)
	http.ListenAndServe(host, nil)
}
