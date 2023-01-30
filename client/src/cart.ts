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

function renderCart(c:Cart) {


}
