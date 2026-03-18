package app

import (
	"encoding/json"
	"examblanc/db"
	"examblanc/models"
	"examblanc/utils"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	judge, err := db.GetJudgeWithName(credentials.Name, true)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not fetch judges from db", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*judge.Password), []byte(credentials.Password))
	if judge == nil || err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(credentials.Name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not generate JWT", http.StatusInternalServerError)
		return
	}

	serialized, err := json.Marshal(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not serialize JWT", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", serialized)
}
