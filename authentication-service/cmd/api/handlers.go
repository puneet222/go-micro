package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	if requestPayload.Email == "" || requestPayload.Password == "" {
		app.errorJSON(w, errors.New("email or password field is required"), http.StatusBadRequest)
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}

	isMatched, err := user.PasswordMatches(requestPayload.Password)

	log.Println("in auth service", isMatched, err)

	if err != nil || !isMatched {
		app.errorJSON(w, errors.New("invalid password"), http.StatusUnauthorized)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("%s is successfully logged in", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, payload)

}
