async function fetchProducts() {
    const response = await fetch('/product-data');
    const products = await response.json();
    return products;
}

function display_product(src: string) {
    var img = document.createElement("img");
    img.src = src;
    return img;
}

async function run() {
    let products = await fetchProducts()
    document.body.appendChild(display_product(products[0].image_path));
        let list = document.createElement("ul");
        for (let x of products) {
            var item = document.createElement("li"); 
            item.innerHTML = x.id + "/" + x.name;
            list.appendChild(item);
        }
        document.getElementById('stuff')!.appendChild(list);
};
run();

type Product = {    name:string, 
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

product = JSON.parse(resp) as Product[]
