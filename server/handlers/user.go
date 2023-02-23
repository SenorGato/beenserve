package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	ApiKey     string `json:"api_key,omitempty"`
	TestApiKey string `json:"test_api_key,omitempty"`
}

type Users struct {
	l *log.Logger
}

func NewUser(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) CreateUser(db_conn *pgx.Conn) func(rw http.ResponseWriter, r *http.Request) {
	if db_conn == nil {
		panic("Nil db_conn in CreateUser")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			panic("Bcrypt hash failed")
		}
		_, err = db_conn.Exec(context.Background(), "INSERT INTO users VALUES($1, $2, $3, $4, $5)",
			user.Email, user.Username, user.Password, user.ApiKey, user.TestApiKey)
		if err != nil {
			return
		}
	}
}
