package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func CookieAuth(db_conn *pgx.Conn) mux.MiddlewareFunc {
	log.Println("Somewhere earlier in the call stack")
	return func(next http.Handler) http.Handler {
		log.Println("Somewhere in the call stack")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			log.Println(r.Cookie("session_id"))
			if err != nil {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
				return
			}
			sessionID := cookie.Value
			if !isValidSession(db_conn, sessionID) {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isValidSession(dbConn *pgx.Conn, sessionId string) bool {
	var session string
	err := dbConn.QueryRow(context.Background(), "SELECT session_id FROM sessions WHERE session_id = $1", sessionId).Scan(&session)
	if err != nil || session == "" {
		log.Println("Error querying database:", err)
		return false
	}
	return true
}
