async function postUserData(reg_form_data:any) {
    fetch('/register', {
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
    body: JSON.stringify(reg_form_data) // body data type must match "Content-Type" header
  })
  console.log(JSON.stringify(reg_form_data))
};

function handleFormSubmit(event: any) {
    event.preventDefault();
    console.log("In handler")
    const data = new FormData(event.target);
    const value = Object.fromEntries(data.entries());
    postUserData(value)
}

function run(){
    const form = document.querySelector('#register');
    form!.addEventListener('submit', handleFormSubmit);
}
run();
