import { renderCreatePost } from "./new-post.js";
import { Logout } from "./authentication.js";

const header = document.getElementById("header");
const main = document.getElementById("main-content");
const sideBar = document.getElementById("sidebar");

function buildHeader() {
	header.innerHTML = `<div class="header-left">
    <h1>Real time forum</h1>
  </div>
  <div class="forum-section">
      <h2 id="welcome-message"></h2>
      <p>Vous êtes connecté au forum.</p>
      <button id="logoutBtn">Déconnexion</button>
    </div>
  <nav class="header-nav">
    <button id="new-post-btn">Nouveau post</button>
    <button id="home-btn">Home</button>
    <button>Catégories</button>
    <button>Chat</button>
  </nav>

  <div class="header-right">
    <span class="user-icon"> </span>
  </div>`;

  document.getElementById("logoutBtn").addEventListener("click", async () => {
    try {
        const response = await fetch("/logout", {
            method: "POST",
            credentials: "include"
        });

        if (response.ok) {
            alert("Déconnecté avec succès !");
        } else {
            alert("Erreur lors de la déconnexion.");
        }
    } catch (err) {
        console.error("Erreur fetch logout :", err);
    }
});


	const postBtn = document.getElementById("new-post-btn");
	postBtn.addEventListener("click", renderCreatePost);
  const logoutBtn = document.getElementById("logoutBtn");
  logoutBtn.addEventListener("click", Logout);
}

function showApp() {
	document.getElementById("auth-container").style.display = "none";
	document.getElementById("app-container").style.display = "block";
	buildHeader();
}

export { header, main, sideBar, buildHeader, showApp };
