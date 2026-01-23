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

  const signUpButton = document.querySelector(".creation-account");
  signUpButton.addEventListener("click", (e) => {
    e.preventDefault();
    // Valider les champs du formulaire ici
    const nicknameInput = document.querySelector("#Nickname");
    const ageInput = document.querySelector("#Age");
    const registerMailInput = document.querySelector("#register-email");
    const genderInput = document.querySelector("#Gender");
    const firstNameInput = document.querySelector("#FirstName");
    const lastNameInput = document.querySelector("#LastName");
    const passwordInput = document.querySelector("#register-password");
    const confirmPasswordInput = document.querySelector(
      "#register-confirm-password",
    );
    let nicknameValue = nicknameInput.value;
    let ageValue = ageInput.value;
    let registerMailValue = registerMailInput.value;
    let genderValue = genderInput.value;
    let lastNameValue = lastNameInput.value;
    let firstNameValue = firstNameInput.value;
    let passwordValue = passwordInput.value;
    let confirmPasswordValue = confirmPasswordInput.value;

    if (passwordValue !== confirmPasswordValue) {
      return "les mots de passe ne correspondent pas";
    } else if (genderValue === "") {
      return "Veuillez sélectionner un genre";
    } else {
      console.log("Nickname: ", nicknameValue);
      console.log("Age : ", ageValue);
      console.log("registerMail : ", registerMailValue);
      console.log("Gender : ", genderValue);
      console.log("lastName : ", lastNameValue);
      console.log("firstName : ", firstNameValue);
      console.log("Mot de passe : ", passwordValue);
      console.log("Confirmation du mot de passe : ", confirmPasswordValue);
    }
  });
}
