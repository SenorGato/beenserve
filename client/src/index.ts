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

function displayProduct(p: Product) {
    console.log(p)
    let product = document.createElement("div")

    let product_name = document.createTextNode(p.name)
    let img = document.createElement("img")
    img.src = p.image_path
    let product_price = document.createTextNode("$" + p.price)
    product.appendChild(product_name)
    product.appendChild(img)
    product.appendChild(product_price)
    document.getElementById("body")!.appendChild(product) || {};
    return product;
}

function addProductToCart(c:Cart, p:Product){
    console.log(`before ${c}`)
    const itemToAdd= c.items.find(item => item.data.sku === p.sku)
    if (itemToAdd) {
        itemToAdd.quantity++
    } else {
        c.items.push({data:p, quantity:1})
    }
    updateCart(c)
    console.log(`after ${c}`)
}

function removeProductFromCart(c:Cart, p:Product){
    console.log(`before remove ${c}`)
    const itemToRemove = c.items.find(item => item.data.sku === p.sku)
    if (itemToRemove) {
        itemToRemove.quantity--
        if (itemToRemove.quantity === 0) {
            c.items = c.items.filter(item => itemToRemove !== item)
        }
    }
    updateCart(c)
    console.log(`after remove ${c}`)
}

function updateCart(c:Cart) {
    console.log("In updateCart")
    console.log(c.items.length)
   const cart_element = document.getElementById("shopping-cart") 
   cart_element!.innerHTML = ""
   for (let i=0; i < c.items.length; i++) {
        const item = c.items[i];
        console.log(`$${item.data.price} ${item.data.name} quant:${item.quantity}`)
   }
}

async function shipCart(c:Cart) {
  const resp = await fetch('/shipcart', {
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
  });
    //const client_key = await resp.json() as string;
    //console.log(client_key)
};

async function run() {
    let shopping_cart: Cart = {items:[]} 
    let product_json = await fetchProducts()
    for (let x in product_json) {
        const product = displayProduct(product_json[x])
        let r_button = document.createElement("button")
        let l_button = document.createElement("button")
        let shipIt = document.createElement("button")
        l_button.onclick = () => {addProductToCart(shopping_cart, product_json[x])}
        r_button.onclick = () => {removeProductFromCart(shopping_cart, product_json[x])}
        shipIt.onclick = () => {shipCart(shopping_cart)}
        product.appendChild(l_button)
        product.appendChild(r_button)
        product.appendChild(shipIt)
        product.appendChild(document.createElement("br"))
    }
};
run();
