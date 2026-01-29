export function ToggleRegisterFrom() {
  const createAccount = document.querySelector(".new-account");
  const registerForm = document.querySelector(".register-form");
  const connectionForm = document.querySelector(".connection-form");

  createAccount.addEventListener("click", (e) => {
    e.preventDefault();
    registerForm.style.display =
      registerForm.style.display === "flex" ? "none" : "flex";
    connectionForm.style.display = "none";
  });
}
