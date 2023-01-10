package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SenorGato/beenserve/src/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v74"
	//_ "github.com/stripe/stripe-go/v72/checkout/session"
)

func main() {
	stripe.Key = "sk_test_51MNgItJUna26uIQEc7yGt2dYnwLjWOrpRSEsnITSK87j3Ff0BB5N7aKs1eOKYwmwEaRNIAnUD7Wz7IWLstq3ovku00vLwGPfEW"
	l := log.New(os.Stdout, "products", log.LstdFlags)

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	sm := mux.NewRouter()
	ph := handlers.NewProducts(l)
	ch := handlers.NewCheckout(l)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts(conn))

	stripeCheckoutRouter := sm.Methods(http.MethodPost).Subrouter()
	stripeCheckoutRouter.HandleFunc("/checkout", ch.CreateCheckoutSession)

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		l.Println("Starting server on port 9090.")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
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
