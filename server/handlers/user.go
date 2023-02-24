package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
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
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		apihash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		testapihash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			panic("Bcrypt hash failed")
		}
		_, err = db_conn.Exec(context.Background(), "INSERT INTO users VALUES($1, $2, $3, $4, $5)",
			user.Email, user.Username, hash, apihash, testapihash)
		if err != nil {
			return
		}
	}
}

func (u *Users) HandleLogin(db_conn *pgx.Conn) func(rw http.ResponseWriter, r *http.Request) {
	if db_conn == nil {
		panic("Nil db_conn in CreateUser")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, err := getHashedPassword(db_conn, loginRequest.Username)
		if err != nil {
			http.Error(rw, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password))
		if err != nil {
			http.Error(rw, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		sessionId, err := generateSessionToken(db_conn, loginRequest.Username)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		writeSessionToDb(sessionId, loginRequest.Username, db_conn)
		http.SetCookie(rw, &http.Cookie{
			Name:     "session_id",
			Value:    sessionId,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		// Return a success response to the front end
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Login successful"))
	}
}

func writeSessionToDb(sessionId string, username string, db_conn *pgx.Conn) {
	_, err := db_conn.Exec(context.Background(),
		"INSERT INTO sessions VALUES($1, $2)", sessionId, username)
	if err != nil {
		panic("DB conn failed in writeSessionToDb")
	}
}

func getHashedPassword(db_conn *pgx.Conn, user string) (string, error) {
	var passHash string
	err := db_conn.QueryRow(context.Background(), "SELECT pass_hash FROM users WHERE username = $1", user).Scan(&passHash)
	return passHash, err
}

func generateSessionToken(db_conn *pgx.Conn, username string) (string, error) {
	// Generate a random byte slice to use as the session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	// Encode the byte slice as a base64 string to create the session token
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return token, nil
}
