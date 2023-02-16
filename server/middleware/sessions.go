package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

func cookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(rw, r, "/error", http.StatusSeeOther)
			return
		}
		sessionID := cookie.Value
		// Check if the session is valid, e.g. by verifying it against a database of valid sessions
		if !isValidSession(sessionID) {
			http.Redirect(rw, r, "/error", http.StatusSeeOther)
			return
		}
		// If the session is valid, pass through to the next handler
		next.ServeHTTP(rw, r)
	})
}

func isValidSession(sessionID string) bool {
	// Query the database or some other data store to see if the session ID is valid
	// Here's a sample implementation using a hard-coded list of valid session IDs
	validSessions := []string{"abc123", "def456", "ghi789"}
	for _, validID := range validSessions {
		if sessionID == validID {
			return true
		}
	}
	return false
}

func HandleLogin(db_conn *pgx.Conn) func(rw http.ResponseWriter, r *http.Request) {
	if db_conn == nil {
		panic("Nil db_conn in CreateUser")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("In handle login")
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		// Query the database to retrieve the hashed password for the user
		hashedPassword, err := getHashedPassword(db_conn, loginRequest.Username, loginRequest.Password)
		if err != nil {
			http.Error(rw, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		// Hash the plaintext password provided in the login request and compare it to the hashed password in the database
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password))
		if err != nil {
			http.Error(rw, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		// If the passwords match, generate a new session token
		sessionID, err := generateSessionToken(db_conn, loginRequest.Username)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// Set the session ID cookie in the response header
		http.SetCookie(rw, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Secure:   true,
		})
		// Return a success response to the front end
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Login successful"))
	}
}

func getHashedPassword(db_conn *pgx.Conn, password string, user string) (string, error) {
	hash, err := db_conn.Query(context.Background(), "SELECT $1 FROM users WHERE username = $2",
		password, user)
	fmt.Println(hash)
	return hash, err
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

// Recieve the password
// Compare the password with the user pass hash in the DB
// If true then issue session token
// Store in cookie and ship to frontend
// Must re-add cookie to all headers sent to backend
