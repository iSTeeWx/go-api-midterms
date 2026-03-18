package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"examblanc/db"
	"examblanc/models"
	"examblanc/utils"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func validateMatch(match models.Match) []string {
	errs := []string{}

	if match.Athlete1Id == nil {
		errs = append(errs, "Field athlete1_id is required")
	} else {
		athlete, err := db.GetAthlete(*match.Athlete1Id)
		if err != nil {
			fmt.Println(err)
			errs = append(errs, "Could not verify athlete1's existence")
		}

		if athlete == nil {
			errs = append(errs, "Athlete1 sohuld exist")
		}
	}

	if match.Athlete2Id == nil {
		errs = append(errs, "Field athlete2_id is required")
	} else {
		athlete, err := db.GetAthlete(*match.Athlete2Id)
		if err != nil {
			fmt.Println(err)
			errs = append(errs, "Could not verify athlete2's existence")
		}

		if athlete == nil {
			errs = append(errs, "Athlete2 sohuld exist")
		}
	}

	if match.JudgeId == nil {
		errs = append(errs, "Field judge_id is required")
	} else {
		judge, err := db.GetJudge(*match.JudgeId)
		if err != nil {
			fmt.Println(err)
			errs = append(errs, "Could not verify judge's existence")
		}
		if judge == nil {
			errs = append(errs, "The judge should exist")
		}
	}

	if match.Date == nil {
		errs = append(errs, "Field date is required")
	} else {
		if time.Unix(int64(*match.Date), 0).After(time.Now()) {
			errs = append(errs, "Match cannot happen in the future")
		}
	}

	if match.Score1 == nil {
		errs = append(errs, "Field score1 is required")
	} else if *match.Score1 < 0 {
		errs = append(errs, "score1 should be positive or null")
	}

	if match.Score2 == nil {
		errs = append(errs, "Field score2 is required")
	} else if *match.Score2 < 0 {
		errs = append(errs, "score2 should be positive or null")
	}

	if match.Athlete1Id != nil &&
		match.Athlete2Id != nil &&
		*match.Athlete1Id == *match.Athlete2Id {
		errs = append(errs, "The two athletes must be different")
	}

	if match.JudgeId != nil && match.Date != nil {
		t := time.Unix(*match.Date, 0)
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		count, err := db.GetCountMatchesByJudgeBetweenDates(*match.JudgeId, start, end)
		if err != nil {
			fmt.Println(err)
			errs = append(errs, "Could not verify that judge isn't overworked")
		}

		if count >= 5 {
			errs = append(errs, "A judge cant oversee more than 5 matches in a day")
		}
	}

	if match.Athlete1Id != nil && match.Athlete2Id != nil {
		count, err := db.GetCountMatchesBetweenAthletesInDay(*match.Athlete1Id, *match.Athlete2Id, time.Unix(*match.Date, 0))
		if err != nil {
			fmt.Println(err)
			errs = append(errs, "Could not verify that athletes are not competing more than once a day")
		}

		if count != 0 {
			errs = append(errs, "Two athletes can't compete more than once a day")
		}
	}

	return errs
}

func PostMatch(w http.ResponseWriter, r *http.Request) {
	var match models.Match

	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("PostMatch(app): empty body")
			http.Error(w, "A json object is required!", http.StatusBadRequest)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to scan request body", http.StatusBadRequest)
		}
		return
	}

	errs := validateMatch(match)
	if len(errs) > 0 {
		serializedErrors, err := json.Marshal(errs)
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
	match.Id = &newUuid

	err = db.PostMatch(match)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not insert match in db", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetMatchesOfJudge(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	tokenString := r.Header.Get("Authorization")
	name, err := utils.VerifyJWT(tokenString)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	judge, err := db.GetJudgeWithName(name, false)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not verify your identity", http.StatusUnauthorized)
		return
	}

	if judge == nil {
		fmt.Println(err)
		http.Error(w, "Token matches no judge", http.StatusUnauthorized)
		return
	}

	if *judge.Id != id {
		fmt.Println(err)
		http.Error(w, "Your token is not that of the target judge", http.StatusUnauthorized)
		return
	}

	matches, err := db.GetMatchesOfJudge(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	serialized, err := json.Marshal(matches)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", serialized)
}


func DeleteMatch(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	err := db.DeleteMatch(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			http.Error(w, "No match with such id", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Error deleting match", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
