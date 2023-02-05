import {loadStripe} from '@stripe/stripe-js';

type Product = {    
    name:string, 
    price:number, 
    sku:string, 
    image_path:string
}

type Cart = {
    items:{data:Product, quantity:number}[]
}

async function fetchProducts() {
    const response = await fetch('/product-data');
    const products = await response.json() as Product[];
    return products;
}

async function fetchPublishableKey() {
        const {publishableKey} = await fetch('/stripe/pubkey').then((r) => r.json());
        if (!publishableKey) {
            alert('Please set your Stripe publishable API key in the .env file');
        } 
        return publishableKey
}

function displayProduct(p: Product):HTMLElement {
    let product_div = document.createElement("div")

    let product_name = document.createTextNode(p.name)
    let img = document.createElement("img")
    img.src = p.image_path
    let product_price = document.createTextNode("$" + p.price)
    product_div.appendChild(product_name)
    product_div.appendChild(img)
    product_div.appendChild(product_price)

    document.getElementById("body")!.appendChild(product_div) || {};
    return product_div;
}

function createButtons(c:Cart, p:Product, div: HTMLElement){
    let r_button = document.createElement("button")
    let l_button = document.createElement("button")
    r_button.textContent = '++++';
    l_button.textContent = '----';
    l_button.onclick = () => {addProductToCart(c, p)}
    r_button.onclick = () => {removeProductFromCart(c, p)}
    div.appendChild(l_button)
    div.appendChild(r_button)
}

function addProductToCart(c:Cart, p:Product){
    const itemToAdd = c.items.find(item => item.data.sku === p.sku)
    if (itemToAdd) {
        itemToAdd.quantity++
    } else {
        c.items.push({data:p, quantity:1})
    }
    updateCart(c)
}

function removeProductFromCart(c:Cart, p:Product){
    const itemToRemove = c.items.find(item => item.data.sku === p.sku)
    if (itemToRemove) {
        itemToRemove.quantity--
        if (itemToRemove.quantity === 0) {
            c.items = c.items.filter(item => itemToRemove !== item)
        }
    }
    updateCart(c)
}

function updateCart(c:Cart) {
    const cart_element = document.getElementById("shopping-cart") as HTMLInputElement
    cart_element.innerHTML = c.items.map(item => `$${item.data.price} ${item.data.name} quant:${item.quantity}`).join("<br>")
    cart_element.innerHTML += "<br>Cart Total: "
    cart_element.innerHTML += c.items.reduce((acc, each) => acc+each.data.price * each.quantity, 0)
}

async function mountPaymentDiv(json: any) {
    var pKey = await fetchPublishableKey()
    const stripe = await loadStripe(pKey);
    const options = {
        clientSecret: json.client_secret,
    };
    const elements = stripe!.elements(options);
    const paymentElement = elements.create('payment');
    paymentElement.mount('#payment-element');
    var e:any
    document.querySelector("#payment-form")!.addEventListener("submit", () => handleSubmit(e, stripe, elements))
}

async function shipCart(c:Cart) {
    fetch('/shipcart', {
    method: 'POST', // *GET, POST, PUT, DELETE, etc.
    mode: 'same-origin', // no-cors, *cors, same-origin
    cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
    credentials: 'include', // include, *same-origin, omit
    headers: {
      'Content-Type': 'application/json'
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
    redirect: 'follow', // manual, *follow, error
    referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(c) // body data type must match "Content-Type" header
  }).then((response) => response.json())
  .then((data) => mountPaymentDiv(data))
};

async function handleSubmit(e:any, stripe:any, ele:any) {
    e.preventDefault();
    //setLoading(true);

    const { error } = await stripe.confirmPayment({
        ele,
        confirmParams: {
            returl_url: "http://localhost:9090/checkout.html"
        },
    });
    //setLoading(false);
}

//function setLoading(isLoading) {
    //if (isLoading) {
        //document.querySelector("#submit")!.disabled = true;
        //document.querySelector("#spinner")!.classList.remove("hidden")
        //document.querySelector("button-text")!.classList.add("hidden")
    //} else {
        //document.querySelector("#submit")!.disabled = false;
        //document.querySelector("#spinner")!.classList.add("hidden")
        //document.querySelector("button-text")!.classList.remove("hidden")
    //}
//}

async function run() {
    let shopping_cart: Cart = {items:[]} 
    let product_json = await fetchProducts()


    let shipIt = document.createElement("button")
    shipIt.textContent = 'LETS GO(SHIP)!';
    shipIt.onclick = () => {shipCart(shopping_cart)}
    document.body.appendChild(shipIt)

    for (let product in product_json) {
        const div: HTMLElement = displayProduct(product_json[product])
        createButtons(shopping_cart, product_json[product], div)
        div.appendChild(document.createElement("br"))
    }
};
run();
