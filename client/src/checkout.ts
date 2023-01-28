//require {loadStripe} from '@stripe/stripe-js';
//import * as Stripe from '@stripe/stripe-js';
import {loadStripe} from '@stripe/stripe-js';
//import {loadStripe} from "https://js.stripe.com/v3/";

const addMessage = (message:string) => {
    const messagesDiv = document.querySelector('#error-message');
    messagesDiv!.innerHTML += ">" + message + '<br>';
    console.log('StripeSampleDebug:', message);
    }

document.addEventListener('DOMContentLoaded', async () => {
    
    const {publishableKey} = await fetch('/stripe/pubkey').then((r) => r.json());
    if (!publishableKey) {
    addMessage(
      'No publishable key returned from the server. Please check `.env` and try again'
    );
    alert('Please set your Stripe publishable API key in the .env file');
    } 
    
    const stripe = await loadStripe('publishableKey');
    const options = {
        clientSecret: document.querySelector('[data-secret]')!.getAttribute('data-secret') || '{}',
    };
    
    const elements = stripe!.elements(options);
    const paymentElement = elements.create('payment');
    paymentElement.mount('#payment-element');

    const form = document.querySelector('#payment-form');
    let submitted = false
    form!.addEventListener('submit', async (e) => {
        console.log("In submit event listener")
        e.preventDefault();

        // Disable double submission of the form
        if(submitted) { return; }
        submitted = true;
        form!.querySelector('button')!.disabled = true;

        // Make a call to the server to create a new
        // payment intent and store its client_secret.
        //const {error: backendError, clientSecret} = await fetch(
        const {} = await fetch(
        '/checkout',
        {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          currency: 'usd',
          paymentMethodType: 'card',
        }),
        }
    ).then((r) => r.json());
});
})
