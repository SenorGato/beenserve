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

//func CreateCheckoutSession() (p *stripe.CheckoutSession.url, err error) {
//params := &stripe.CheckoutSessionParams{
//Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
//LineItems: []*stripe.CheckoutSessionLineItemParams{
//{
//PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
//Currency: stripe.String("usd"),
//ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
//Name: stripe.String("T-shirt"),
//},
//UnitAmount: stripe.Int64(2000),
//},
//Quantity: stripe.Int64(1),
//},
//},
//SuccessURL: stripe.String("http://localhost:9090/success"),
//CancelURL:  stripe.String("http://localhost:9090/cancel"),
//}
//s, _ := session.New(params)
//
////if err != nil {
////return s
////}
//
//return s
//}
