async function fetchProducts() {
    const response = await fetch('/product-data');
    const products = await response.json();
    return products;
}

async function run() {
    let products = await fetchProducts()
    console.log(products);
        let list = document.createElement("ul");
        for (let x of products) {
            var item = document.createElement("li"); 
            item.innerHTML = x.id + "/" + x.name;
            list.appendChild(item);
        }
        document.getElementById('stuff').appendChild(list);
};
run();
