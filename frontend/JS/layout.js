import {renderCreatePost} from "./new-post.js";
import {Logout} from "./authentication.js";
import {handleChatClick} from "./chat.js";

const header = document.getElementById("header");
const main = document.getElementById("main-content");
const sideBar = document.getElementById("sidebar");

function buildHeader() {
  header.innerHTML = `<div class="header-left">
    <h1>Real time forum</h1>
  </div>
  <nav class="header-nav">
    <button id="new-post-btn">Nouveau post</button>
    <button id="home-btn">Home</button>
    <button id="categories-btn">Catégories</button>
    <button id="chat-btn">Chat</button>
  </nav>
  <div class="forum-section">
  <div class="profile-section">
      <p id="welcome-message"></p>
      <img src="./frontend/img/profil.gif" alt="image profil" class="profil-icon">
      </div>
      <button id="logoutBtn">Déconnexion</button>
    </div>
`;

  document.getElementById("logoutBtn").addEventListener("click", async () => {
    try {
      const response = await fetch("/logout", {
        method: "POST",
        credentials: "include",
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
  const chatBtn = document.querySelector("#chat-btn");
  chatBtn ? console.log("bouton exist") : console.log("bouton introuvable");
  chatBtn.addEventListener("click", (e) => {
    console.log("click chat");
    handleChatClick(e);
  });
}
function buildSidebar() {
  sideBar.innerHTML = `<h2>Users</h2>
  <div class="users-list"></div>
  `;
}

function buildMain() {
  main.innerHTML = `<h2>Posts</h2>
  <div class="posts-header">
  			<span>Titre</span>
  			<span>Catégorie(s)</span>
  			<span>Texte</span>
		</div>`;
}

function showApp() {
  document.getElementById("auth-container").style.display = "none";
  document.getElementById("app-container").style.display = "contents";
  buildHeader();
  buildSidebar();
  buildMain();
}

export {header, main, sideBar, buildHeader, showApp};
