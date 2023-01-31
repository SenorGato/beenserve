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
	"github.com/joho/godotenv"
	"github.com/senorgato/beenserve/server/handlers"
	"github.com/stripe/stripe-go/v74"
)

func main() {
	err := godotenv.Load("test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripe.Key = os.Getenv("STRIPE_KEY")
	addr := os.Getenv("WEB_SERVER_PORT")

	product_log := log.New(os.Stdout, "Products:", log.LstdFlags)
	checkout_log := log.New(os.Stdout, "Checkout:", log.LstdFlags)
	server_log := log.New(os.Stdout, "Server:", log.LstdFlags)

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	sm := mux.NewRouter()

	ph := handlers.NewProducts(product_log)
	ch := handlers.NewCheckout(checkout_log)

	// Database route
	productRouter := sm.Methods(http.MethodGet).Subrouter()
	productRouter.HandleFunc("/product-data", ph.GetProducts(conn))
	// Stripe-API routes
	stripeCheckoutRouter := sm.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	stripeCheckoutRouter.HandleFunc("/checkout", ch.CreateCheckoutSession).Methods("GET", "POST")
	stripeCheckoutRouter.HandleFunc("/stripe/pubkey", ch.PubKey).Methods("GET", "POST")
	// stripeCheckoutRouter.HandleFunc("/shipcart", ch.RecieveCart)

	// Static Files
	fs := http.FileServer(http.Dir("/go/bin/client"))
	sm.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	s := http.Server{
		Addr:         ":" + addr,
		Handler:      sm,
		ErrorLog:     server_log,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		server_log.Println("Starting server on port 9090.")

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
