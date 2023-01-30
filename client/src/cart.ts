//product = JSON.parse(resp) as Product[]

type Product = {    
    name:string, 
    price:number, 
    sku:string, 
    path:string
}
type Cart = {
    items:{data:Product, quantity:number}[]
}

function addProductToCart(c:Cart, p:Product){
    for (const item of c.items) {
        if (item.data.sku === p.sku) {
            item.quantity++
            return
        }
    }
    c.items.push({data:p, quantity:1})
}

function removeProductFromCart(c:Cart, p:Product){
    const itemToRemove = c.items.find(item => item.data.sku === p.sku);
    if (itemToRemove) {
        itemToRemove.quantity--
        if (itemToRemove.quantity === 0) {
            c.items = c.items.filter(item => itemToRemove !== item)
        }
    }
}

async function shipCart(c:Cart) {
  await fetch('/shipcart', {
    method: 'POST', // *GET, POST, PUT, DELETE, etc.
    mode: 'cors', // no-cors, *cors, same-origin
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
};
