async function postLoginData(login_form_data:any) {
    fetch('/login', {
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
    body: JSON.stringify(login_form_data) // body data type must match "Content-Type" header
  })
};

async function handleLoginSubmit(event: any) {
    event.preventDefault();
    const formData = new FormData(event.target);
    let value = Object.fromEntries(formData.entries());
    postLoginData(value);
}

function login_run(){
    const form = document.querySelector('#login');
    form!.addEventListener('submit', handleLoginSubmit);
}

login_run();
