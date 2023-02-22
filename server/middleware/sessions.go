package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash struct {
	hash string
}

func HandleLogin(db_conn *pgx.Conn) func(rw http.ResponseWriter, r *http.Request) {
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
		hashedPassword, err := getHashedPassword(db_conn, loginRequest.Username, loginRequest.Password)
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
		})
		// Return a success response to the front end
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Login successful"))
	}
}

func writeSessionToDb(sessionId string, username string, db_conn *pgx.Conn) {
	_, err := db_conn.Exec(context.Background(), "INSERT INTO sessions (session_key, user) VALUES($1, $2)", sessionId, username)
	if err != nil {
		panic("DB conn failed in writeSessionToDb")
	}
}

func getHashedPassword(db_conn *pgx.Conn, password string, user string) (string, error) {
	rows, err := db_conn.Query(context.Background(), "SELECT $1 FROM users WHERE username = $2",
		password, user)
	defer rows.Close()
	results, err := pgx.CollectRows(rows, pgx.RowToStructByPos[PasswordHash])
	return results[0].hash, err
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

// Recieve the password
// Compare the password with the user pass hash in the DB
// If true then issue session token
// Store in cookie and ship to frontend
// Must re-add cookie to all headers sent to backend
