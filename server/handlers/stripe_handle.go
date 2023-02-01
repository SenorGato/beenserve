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

func NewCheckout(l *log.Logger) *Checkout {
	return &Checkout{l}
}

type CheckoutData struct {
	ClientSecret string
}

type paymentIntentCreateReq struct {
	Currency          string `json:"currency"`
	PaymentMethodType string `json:"paymentMethodType"`
}

type Cart struct {
	Data     string
	Quantity int
}

//func (c *Checkout) calculateTotal(cart Cart) (cost int) {
//for i, v := range cart{
//c.l.Println("Index:%d Value:%+v", i, v)
//}
//return cost;
//}

func (c *Checkout) RecieveCart(rw http.ResponseWriter, r *http.Request) {
	var cart Cart
	c.l.Println("The cart has been recieved.")
	err := json.NewDecoder(r.Body).Decode(&cart)
	c.l.Println(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	c.l.Println(os.Getenv("STRIPE_KEY"))
	checkoutTmpl, err := template.ParseFiles("./client/html/checkout.html")
	if err != nil {
		panic(err)
	}

	stripe.Key = os.Getenv("STRIPE_KEY")
	c.l.Println(stripe.Key)
	// req := paymentIntentCreateReq{}
	// json.NewDecoder(r.Body).Decode(&req)

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
	// fmt.Fprintf(rw, "Cart: %+v", cart)
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
