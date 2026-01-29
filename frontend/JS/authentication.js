export function Authentication() {
    const connectionForm = document.querySelector(".connection-form")

    connectionForm.addEventListener("submit", (e) => {
        e.preventDefault();

        const username = document.getElementById("register-username").value;
        const password = document.getElementById("password").value;

        fetch("/auth/login", { // envoie la requête au serveur sur la route /auth/login
            method: "POST", // on envoie des données
            headers: {
                "Content-Type": "application/json" // le serveur sait qu’on envoie du JSON
            },
            body: JSON.stringify({ // on transforme nos valeurs JS en JSON
                username: username,
                password: password
            })
        })
        .then(response => response.json()) // le serveur renvoie du JSON, on le transforme en objet JS
        .then(data => { // Gestion de la réponse (succès ou erreur)
            console.log("Réponse serveur :", data);
            if (data.success) {
                // Masquer les formulaires
                connectionForm.style.display = "none";
                document.querySelector(".register-form").style.display = "none";

                // Afficher la section forum
                const forumSection = document.querySelector(".forum-section");
                forumSection.style.display = "block";

                // Mettre le pseudo dans le message
                document.getElementById("welcome-message").textContent = `Bienvenue, ${data.user.nickname} !`;

                // Gérer le logout
                document.getElementById("logout").addEventListener("click", () => {
                    forumSection.style.display = "none";
                    connectionForm.style.display = "flex";
                    document.querySelector(".register-form").style.display = "none";
                });
            } else {
                alert("Connexion échouée !");
            }
        })
        .catch(error => { // attrape les erreurs réseau
            console.log("Erreur réseau :", error)
        });

        console.log("Username :", username, "Password :", password)
    })
}