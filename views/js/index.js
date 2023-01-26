async function fetchProducts() {
    const response = await fetch('/product-data');
    const products = await response.json();
    return products;
}

function display_product(src) {
    var img = document.createElement("img");
    img.src = src;
    document.body.appendChild(img);
}

async function run() {
    let products = await fetchProducts()
    show_image(products[0].image_path);
    //console.log(products);
    //console.log(products[0]);
        let list = document.createElement("ul");
        for (let x of products) {
            var item = document.createElement("li"); 
            item.innerHTML = x.id + "/" + x.name;
            list.appendChild(item);
        }
        document.getElementById('stuff').appendChild(list);
};
run();
