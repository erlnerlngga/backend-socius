package util

import (
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func SendEmail(email, token string) error {
	// set email auth
	authEmail := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("PASSWORD_EMAIL"), "smtp.gmail.com")

	// compose email
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" + "Subject: Sign In Link\r\n" + "\r\n" + "http://localhost:3000/auth/" + token)

	if err := smtp.SendMail("smtp:gmail.com:587", authEmail, "laann.en@gmail.com", to, msg); err != nil {
		return err
	}

	return nil
}
