import bcrypt from 'bcryptjs';

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
};

async function hashRegisterForm(reg_form_data: any) {
    const saltRounds = 10;
    const pass_salt = await bcrypt.genSalt(saltRounds);
    const api_salt= await bcrypt.genSalt(saltRounds);
    const test_api_salt = await bcrypt.genSalt(saltRounds);

    reg_form_data.password = await bcrypt.hash(reg_form_data.password, pass_salt);
    reg_form_data.api_key = await bcrypt.hash(reg_form_data.api_key, api_salt);
    reg_form_data.test_api_key = await bcrypt.hash(reg_form_data.test_api_key, test_api_salt);

    return reg_form_data;
}

async function handleFormSubmit(event: any) {
    event.preventDefault();
    const data = new FormData(event.target);
    const value = Object.fromEntries(data.entries());
    const hashedForm = await hashRegisterForm(value);
    postUserData(hashedForm)
}

function run(){
    const form = document.querySelector('#register');
    form!.addEventListener('submit', handleFormSubmit);
}
run();
