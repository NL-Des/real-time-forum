export function CreateAccount() {
  const comeBack = document.querySelector(".new-account-comeback");
  const registerForm = document.querySelector(".register-form");
  const connectionForm = document.querySelector(".connection-form");

  comeBack.addEventListener("click", (e) => {
    e.preventDefault();
    // Masquer le formulaire d'inscription
    registerForm.style.display = "none";

    // Afficher le formulaire de connexion
    connectionForm.style.display = "flex";
  });
}
