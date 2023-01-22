package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type Checkout struct {
	l *log.Logger
}

type CheckoutData struct {
	ClientSecret string
}

func NewCheckout(l *log.Logger) *Checkout {
	return &Checkout{l}
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

type paymentIntentCreateReq struct {
	Currency          string `json:"currency"`
	PaymentMethodType string `json:"paymentMethodType"`
}

func (c *Checkout) CreateCheckoutSession(rw http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	checkoutTmpl, err := template.ParseFiles("/home/senoraraton/bins/beenserve/views/checkout.html")
	if err != nil {
		panic(err)
	}
	req := paymentIntentCreateReq{}
	fmt.Printf("request created:%s", req)
	json.NewDecoder(r.Body).Decode(&req)
	fmt.Printf("request created:%s", req)

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(5999),
		Currency:           stripe.String("usd"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}
	intent, err := paymentintent.New(params)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			fmt.Printf("Other Stripe error occurred: %v\n", stripeErr.Error())
		} else {
			fmt.Printf("Other error occurred: %v\n", err.Error())
		}
		return
	}
	data := CheckoutData{
		ClientSecret: intent.ClientSecret,
	}
	checkoutTmpl.Execute(rw, data)
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
