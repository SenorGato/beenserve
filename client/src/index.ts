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
