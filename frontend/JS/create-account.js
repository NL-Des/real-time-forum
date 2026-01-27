export function CreateAccount() {
  const comeBack = document.querySelector(".new-account-comeback");
  const registerForm = document.querySelector(".register-form");
  const connectionForm = document.querySelector(".connection-form");
  const successMessage = document.getElementById("success-message");
  const errorMessage = document.getElementById("error-message");
  const registerFormElement = document.getElementById("register-form");
  const registerFormPassword = document.getElementById("register-password");
  const registerFormConfirmPassword = document.getElementById(
    "register-confirm-password",
  );

  // Cacher les messages au départ
  successMessage.style.display = "none";
  errorMessage.style.display = "none";

  // Bouton retour vers connexion
  comeBack.addEventListener("click", (e) => {
    e.preventDefault();
    registerForm.style.display = "none";
    connectionForm.style.display = "flex";
  });

  // Gestion de la soumission du formulaire
  registerFormElement.addEventListener("submit", async function (e) {
    e.preventDefault();

    // Cacher les messages précédents
    successMessage.style.display = "none";
    errorMessage.style.display = "none";

    // Récupérer les données du formulaire
    const formData = new FormData(this);

    try {
      const response = await fetch("/", {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        successMessage.textContent = "Account created successfully!";
        successMessage.style.display = "block";

        // Réinitialiser le formulaire
        registerFormElement.reset();

        // Optionnel : Basculer vers le formulaire de connexion après 2 secondes
        setTimeout(() => {
          registerForm.style.display = "none";
          connectionForm.style.display = "flex";
          successMessage.style.display = "none";
        }, 2000);
      } else {
        const errorText = await response.text();
        errorMessage.textContent = errorText;
        errorMessage.style.display = "block";
        registerFormPassword.reset();
        registerFormConfirmPassword.reset();

        // Faire défiler vers le haut pour voir l'erreur
        window.scrollTo({top: 0, behavior: "smooth"});
      }
    } catch (error) {
      errorMessage.textContent = "Network error. Please try again.";
      errorMessage.style.display = "block";
    }
  });
}
