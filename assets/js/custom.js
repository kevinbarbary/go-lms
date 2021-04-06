function signIn() {
    let none = document.getElementById("sign-in-none");
    none.classList.add("spinner-border");
    none.classList.add("spinner-border-sm");
    let btn = document.getElementById("sign-in-btn");
    btn.classList.add("disabled");
    btn.classList.add("btn-text");
    let text = document.getElementById("sign-in-text");
    text.classList.add("visually-hidden");
}
