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

type Cart struct {
	Items []Item `json:"items,omitempty"`
}
type Item struct {
	Data     Prod    `json:"data,omitempty"`
	Quantity float32 `json:"quantity,omitempty"`
}
type Prod struct {
	Price      float32 `json:"price,omitempty"`
	Sku        string  `json:"sku,omitempty"`
	Name       string  `json:"name,omitempty"`
	Image_path string  `json:"image_path,omitempty"`
}

type Tester struct {
	Name []Test_Two `json:"name"`
}

type Test_Two struct {
	Data     Prod    `json:"data,omitempty"`
	Quantity float32 `json:"quantity,omitempty"`
}

func (c *Checkout) RecieveCart(rw http.ResponseWriter, r *http.Request) {
	var cart Cart
	// var test Tester
	c.l.Println("The cart has been recieved.")
	// body, err := io.ReadAll(r.Body)
	// myString := string(body[:])
	// c.l.Println(myString)
	// c.l.Println(cart.Quantity)
	// json_body := json.Unmarshal(body, &cart)
	// c.l.Println(json_body)
	// c.l.Println(json.Unmarshal(body, &cart))
	// c.l.Println(cart.Data.Price)
	err := json.NewDecoder(r.Body).Decode(&cart)
	c.l.Println(cart)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	checkoutTmpl, err := template.ParseFiles("./client/html/checkout.html")
	if err != nil {
		panic(err)
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

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

//func (c *Checkout) calculateTotal(cart Cart) (cost int) {
//for i, v := range cart{
//c.l.Println("Index:%d Value:%+v", i, v)
//}
//return cost;
//}
