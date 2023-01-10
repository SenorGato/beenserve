// Set your publishable key: remember to change this to your live publishable key in production
// See your keys here: https://dashboard.stripe.com/apikeys
const stripe = stripe('pk_test_51MNgItJUna26uIQEAZhCYdVAvrc0pM7qtJvKP6oe8lEbgcefGL9hEhLeoOZYaxklq0ih6enZVnwMw8DzO2VY5Tmj00njjYCafM');

const options = {
  clientSecret: '{{CLIENT_SECRET}}',
  // Fully customizable with appearance API.
  appearance: {/*...*/},
};

// Set up Stripe.js and Elements to use in checkout form, passing the client secret obtained in step 3
const elements = stripe.elements(options);

// Create and mount the Payment Element
const paymentElement = elements.create('payment');
paymentElement.mount('#payment-element');
