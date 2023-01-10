package handlers

import (
	"fmt"
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

func (c *Checkout) CreateCheckoutSession(rw http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("TESTSTRIPE_KEY")
	checkoutTmpl, err := template.ParseFiles("/home/senoraraton/bins/beenserve/views/checkout.html")
	if err != nil {
		panic(err)
	}
	fmt.Println("In Create Checkout")
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1000),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	intent, _ := paymentintent.New(params)
	data := CheckoutData{
		ClientSecret: intent.ClientSecret,
	}
	checkoutTmpl.Execute(rw, data)
}
