document.addEventListener('DOMContentLoaded', async () => {
    const stripe = Stripe('pk_test_51MNgItJUna26uIQEAZhCYdVAvrc0pM7qtJvKP6oe8lEbgcefGL9hEhLeoOZYaxklq0ih6enZVnwMw8DzO2VY5Tmj00njjYCafM');
    
    const options = {
        clientSecret: document.querySelector('[data-secret]').getAttribute('data-secret'),
        appearance: {/*...*/},
    };
    const elements = stripe.elements(options);
    const paymentElement = elements.create('payment');
    paymentElement.mount('#payment-element');

    const form = document.querySelector('#payment-form');
    form.addEventListener('submit')
});

const addMessage = (message) => {
    const messagesDiv = document.querySelector('#error-message');
    messagesDiv.style.display = 'block';
    messagesDiv.innerHtml += ">" + message + '<br>';
    console.log('StripeSampleDebug:', message);
}
