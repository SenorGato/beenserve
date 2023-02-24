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
	"github.com/senorgato/beenserve/server/handlers"
	"github.com/stripe/stripe-go/v74"
)

func main() {
	stripe.Key = os.Getenv("STRIPE_KEY")
	addr := os.Getenv("WEB_SERVER_PORT")

	product_log := log.New(os.Stdout, "Products:", log.LstdFlags)
	checkout_log := log.New(os.Stdout, "Checkout:", log.LstdFlags)
	server_log := log.New(os.Stdout, "Server:", log.LstdFlags)
	user_log := log.New(os.Stdout, "User:", log.LstdFlags)

	db_conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	sm := mux.NewRouter()

	ph := handlers.NewProducts(product_log)
	ch := handlers.NewCheckout(checkout_log)

	// Database route
	productRouter := sm.Methods(http.MethodGet).Subrouter()
	productRouter.HandleFunc("/product-data", ph.GetProducts(db_conn))
	// Stripe-API routes
	stripeCheckoutRouter := sm.Methods(http.MethodGet, http.MethodPost).Subrouter()
	stripeCheckoutRouter.HandleFunc("/stripe/pubkey", ch.PubKey).Methods("GET", "POST")
	stripeCheckoutRouter.HandleFunc("/shipcart", ch.RecieveCart(db_conn)).Methods("POST")

	// Static Files
	fs := http.FileServer(http.Dir("/go/bin/client"))
	sm.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	userHandler := handlers.NewUser(user_log)

	userRouter := sm.Methods(http.MethodPost).Subrouter()
	userRouter.HandleFunc("/register", userHandler.CreateUser(db_conn))
	userRouter.HandleFunc("/login", userHandler.HandleLogin(db_conn))

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You have access to the protected resource."))
	})
	protectedRoute := sm.Path("/protected").Subrouter()
	protectedRoute.Use(middleware.CookieAuth(db_conn))
	protectedRoute.Methods(http.MethodGet).Handler(protectedHandler)

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
