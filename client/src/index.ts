type RegistrationForm = {
    Username:string,
    Email:string,
    Password:string,
    StripeAPIkey:string,
    TestStripeAPIkey:string
}
function run(){

    const form: HTMLFormElement = document.querySelector('#register')!;
    form.onsubmit = () => {
        const formData = new FormData(form);
        console.log(formData)
        const text = formData.get('textInput') as string;
        console.log(text);
        return false; // prevent reload
    }
}
run();


