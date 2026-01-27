export function Authentication() {
    const connectionForm = document.querySelector(".connection-form").value

    connectionForm.addEventListener("submit", (e) => {
        e.preventDefault();

        const username = document.getElementById("register-username").value;
        const password = document.getElementById("password").value;

        console.log("Username :", username, "Password :", password)
    })
}