package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/senorgato/beenserve-admin/server/handlers"
	"github.com/senorgato/beenserve-admin/server/middleware"
)

func main() {
	server_log := log.New(os.Stdout, "Server:", log.LstdFlags)
	user_log := log.New(os.Stdout, "User:", log.LstdFlags)

	db_conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	sm := mux.NewRouter()

	user_handler := handlers.NewUser(user_log)

	userRouter := sm.Methods(http.MethodPost).Subrouter()
	userRouter.HandleFunc("/register", user_handler.CreateUser(db_conn))
	userRouter.HandleFunc("/login", middleware.HandleLogin(db_conn))

	fs := http.FileServer(http.Dir("/go/bin/client"))
	sm.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	s := http.Server{
		Addr:         ":9091",
		Handler:      sm,
		ErrorLog:     server_log,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		server_log.Println("Starting server on port 9091.")

		err := s.ListenAndServe()
		if err != nil {
			server_log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
