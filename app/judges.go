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
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func GetJudges(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	var pageInt *int
	if page != "" {
		v, err := strconv.Atoi(page)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to get value of page", http.StatusBadRequest)
			return
		}
		if v < 1 {
			fmt.Println(err)
			http.Error(w, "page should be greater than 1", http.StatusBadRequest)
			return
		}
		pageInt = &v
	}

	judges, err := db.GetJudges(pageInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch judges", http.StatusInternalServerError)
		return
	}

	serialized, err := json.Marshal(judges)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing judges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", serialized)
}

func validateJudge(judge models.Judge) []string {
	errors := []string{}

	if judge.Name == nil {
		errors = append(errors, "Field \"name\" is required")
	} else {
		if len(*judge.Name) < 2 {
			errors = append(errors, "name should be longer than 1 characters")
		}
		if len(*judge.Name) > 50 {
			errors = append(errors, "name should be shorter than 51 characters")
		}
		nameError, err := regexp.Match("[^a-zA-Z]", []byte(*judge.Name))
		if err != nil {
			errors = append(errors, "THIS SHOULD NEVER HAPPEN (hopefully)")
		}
		if nameError {
			errors = append(errors, "name should only contain letters")
		}
	}

	if judge.Password == nil {
		errors = append(errors, "Field \"password\" is required")
	} else {
		if len(*judge.Password) < 6 {
			errors = append(errors, "password should be longer than 5 characters")
		}
		if len(*judge.Password) > 50 {
			errors = append(errors, "password should be shorter than 51 characters")
		}
		specialCharsMatch, err := regexp.Match("[!+*/]", []byte(*judge.Password))
		if err != nil {
			errors = append(errors, "THIS SHOULD NEVER HAPPEN (hopefully)")
		}
		if !specialCharsMatch {
			errors = append(errors, "password should contain a special character (!+*/)")
		}
	}

	if judge.ExperienceYears != nil && *judge.ExperienceYears < 0 {
		errors = append(errors, "experience_years should be at least 0")
	}

	if judge.Phone == nil {
		errors = append(errors, "Field \"phone\" is required")
	} else {
		if len(*judge.Phone) != 10 {
			errors = append(errors, "phone number should be exactly 10 characters long")
		}
		if !strings.HasPrefix(*judge.Phone, "0") {
			errors = append(errors, "phone number should start with '0'")
		}
		phoneError, err := regexp.Match("[^0-9]", []byte(*judge.Phone))
		if err != nil {
			errors = append(errors, "THIS SHOULD NEVER HAPPEN (hopefully)")
		}
		if phoneError {
			errors = append(errors, "phone number should only contain numbers")
		}
	}

	return errors
}

func PostJudge(w http.ResponseWriter, r *http.Request) {
	var judge models.Judge

	err := json.NewDecoder(r.Body).Decode(&judge)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("PostJusdge(app): empty body")
			http.Error(w, "A json object is required!", http.StatusBadRequest)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to scan request body", http.StatusBadRequest)
		}
		return
	}

	errs := validateJudge(judge)
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

	exist, err := db.GetJudgeWithName(*judge.Name, false)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to verify unicity of judges name", http.StatusInternalServerError)
		return
	}

	if exist != nil {
		fmt.Println("A judge with this name already exists")
		http.Error(w, "A judge with this name already exists", http.StatusBadRequest)
		return
	}

	newUuid := uuid.New().String()
	judge.Id = &newUuid

	hashed, _ := bcrypt.GenerateFromPassword([]byte(*judge.Password),
		bcrypt.DefaultCost)
	stringHash := string(hashed)
	judge.Password = &stringHash

	err = db.PostJudge(judge)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add athlete to db", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetJudge(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	judge, err := db.GetJudge(id)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch judges", http.StatusInternalServerError)
		return
	}

	if judge == nil {
		fmt.Println("No judge with such id")
		http.Error(w, "No judge with such id", http.StatusNotFound)
		return
	}

	serialized, err := json.Marshal(judge)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error serializing judge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", serialized)
}

func DeleteJudge(w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteJudge(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			http.Error(w, "No judge with such id", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Error deleting judge", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
