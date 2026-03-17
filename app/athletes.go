package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"examblanc/db"
	"examblanc/models"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GetAthletes(w http.ResponseWriter, r *http.Request) {
	athletes, err := db.GetAthletes()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching athletes", http.StatusInternalServerError)
		return
	}

	filterCountry := strings.ToLower(r.URL.Query().Get("country"))
	if filterCountry != "" {
		filteredAthletes := []models.Athlete{}

		for _, athlete := range athletes {
			athleteCountry := strings.ToLower(*athlete.Country)

			if strings.Contains(filterCountry, athleteCountry) ||
				strings.Contains(athleteCountry, filterCountry) {
				filteredAthletes = append(filteredAthletes, athlete)
			}
		}

		athletes = filteredAthletes
	}

	serialized, err := json.Marshal(athletes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing athletes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", serialized)
}

func validateAthlete(athlete models.Athlete) []string {
	errors := []string{}

	if athlete.Name == nil {
		errors = append(errors, "Field \"name\" is required")
	} else {
		if len(*athlete.Name) < 2 {
			errors = append(errors, "name should be longer than 1 characters")
		}
		if len(*athlete.Name) > 50 {
			errors = append(errors, "name should be shorter than 51 characters")
		}
		nameError, err := regexp.Match("[^a-zA-Z]", []byte(*athlete.Name))
		if err != nil {
			errors = append(errors, "THIS SHOULD NEVER HAPPEN (hopefully)")
		}
		if nameError {
			errors = append(errors, "name should only contain letters")
		}
	}

	if athlete.Country == nil {
		errors = append(errors, "Field \"country\" is required")
	} else {
		if len(*athlete.Country) < 3 {
			errors = append(errors, "country should be longer than 2 characters")
		}
		if len(*athlete.Country) > 50 {
			errors = append(errors, "country should be shorter than 51 characters")
		}
	}

	if athlete.Age == nil {
		errors = append(errors, "Field \"age\" is required")
	} else {
		if *athlete.Age < 12 || *athlete.Age > 60 {
			errors = append(errors, "age should be between 12 and 60")
		}
	}

	return errors
}

func PostAthletes(w http.ResponseWriter, r *http.Request) {
	var athlete models.Athlete

	err := json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("PostAtheletes(app): empty body")
			http.Error(w, "A json object is required!", http.StatusBadRequest)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to scan request body", http.StatusBadRequest)
		}
		return
	}

	errors := validateAthlete(athlete)

	if len(errors) > 0 {
		serializedErrors, err := json.Marshal(errors)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "There were errors, but we failed to serialize them for you :(", http.StatusInternalServerError)
			return
		}

		fmt.Println(string(serializedErrors))
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(serializedErrors), http.StatusBadRequest)
		return
	}

	newUuid := uuid.New().String()
	athlete.Id = &newUuid

	err = db.AddAthlete(athlete)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add athlete to db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetAthlete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	athlete, err := db.GetAthlete(id)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching athlete", http.StatusInternalServerError)
		return
	}

	if athlete == nil {
		fmt.Println("No athlete with such id")
		http.Error(w, "No athlete with such id", http.StatusNotFound)
		return
	}

	serialized, err := json.Marshal(athlete)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing athlete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", serialized)
}

func PutAthlete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	exists, err := db.GetAthlete(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching athlete", http.StatusInternalServerError)
		return
	}

	if exists == nil {
		fmt.Println("No athlete with such id")
		http.Error(w, "No athlete with such id", http.StatusNotFound)
		return
	}

	var athlete models.Athlete

	err = json.NewDecoder(r.Body).Decode(&athlete)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("PostAtheletes(app): empty body")
			http.Error(w, "A json object is required!", http.StatusBadRequest)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to scan request body", http.StatusBadRequest)
		}
		return
	}

	errors := validateAthlete(athlete)

	if len(errors) > 0 {
		serializedErrors, err := json.Marshal(errors)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "There were errors, but we failed to serialize them for you :(", http.StatusInternalServerError)
			return
		}

		fmt.Println(string(serializedErrors))
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(serializedErrors), http.StatusBadRequest)
		return
	}

	athlete.Id = exists.Id

	err = db.PutAthlete(*athlete.Id, athlete)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add athlete to db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAthlete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	err := db.DeleteAthlete(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			http.Error(w, "No athlete with such id", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Error deleting athlete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
