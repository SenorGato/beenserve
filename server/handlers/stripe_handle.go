package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type Checkout struct {
	l *log.Logger
}

func NewCheckout(l *log.Logger) *Checkout {
	return &Checkout{l}
}

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

type Cart struct {
	Items []Item `json:"items,omitempty"`
}

type Item struct {
	Data     Product `json:"data,omitempty"`
	Quantity int64   `json:"quantity,omitempty"`
}

func (c *Checkout) CalculateTotal(db_conn *pgx.Conn, shopping_cart Cart) (total int64) {
	if db_conn == nil {
		panic("Nil db_conn in CalculateTotal!")
	}
	for _, s := range shopping_cart.Items {
		rows, err := db_conn.Query(context.Background(), "SELECT price from frontend where SKU = ($1)", s.Data.SKU)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
			os.Exit(1)
		}

		price, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
		total += price * s.Quantity
		if err != nil {
			fmt.Printf("Collect rows error: %v", err)
		}
	}
	return total
}

func (c *Checkout) RecieveCart(db_conn *pgx.Conn) func(rw http.ResponseWriter, r *http.Request) {
	if db_conn == nil {
		panic("Nil db_conn in RecieveCart!")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		stripe.Key = os.Getenv("STRIPE_KEY")
		var cart Cart
		err := json.NewDecoder(r.Body).Decode(&cart)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		cartTotal := c.CalculateTotal(db_conn, cart)
		params := &stripe.PaymentIntentParams{
			Amount:             stripe.Int64(cartTotal * 100),
			Currency:           stripe.String("usd"),
			PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		}

		intent, err := paymentintent.New(params)
		// c.l.Println(intent)
		data := CheckoutData{
			ClientSecret: intent.ClientSecret,
		}
		rw.Header().Set("Content-Type", "application/json")
		// rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(data)
		if err != nil {
			if stripeErr, ok := err.(*stripe.Error); ok {
				fmt.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
			} else {
				fmt.Printf("Other error occurred: %v\n", err.Error())
			}
			return
		}
	}
}

func (c *Checkout) PubKey(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	writeJSON(rw, struct {
		PublishableKey string `json:"publishableKey"`
	}{
		PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

//func (c *Checkout) calculateTotal(cart Cart) (cost int) {
//for i, v := range cart{
//c.l.Println("Index:%d Value:%+v", i, v)
//}
//return cost;
//}
