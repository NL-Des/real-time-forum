import {renderCreatePost} from "./new-post.js";
import {Logout} from "./authentication.js";
import {handleChatClick} from "./chat.js";
import {initWebSocket, getWebSocket} from "./websocket.js"; // ✅ Import

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
  </div>`;

  document
    .getElementById("new-post-btn")
    .addEventListener("click", renderCreatePost);
  document.getElementById("logoutBtn").addEventListener("click", Logout);
  document
    .getElementById("chat-btn")
    .addEventListener("click", handleChatClick);
}

function buildSidebar() {
  sideBar.innerHTML = `<h2>Utilisateurs en ligne</h2>
  <div class="users-list"></div>`;

  // ✅ Initialiser WebSocket et écouter les utilisateurs
  const ws = initWebSocket();

  ws.onmessage = function (event) {
    const data = JSON.parse(event.data);

    // ✅ Mise à jour de la liste des utilisateurs
    if (data.type === "online_users") {
      updateUsersList(data.users);
    }

    // ✅ Notification pour un nouveau message
    if (data.type === "message") {
      const userItem = document.querySelector(
        `.user-item[data-user-id="${data.sender_id}"]`,
      );
      if (userItem) {
        userItem.classList.add("has-notification");
      }
    }
  };
}

// ✅ Fonction pour mettre à jour la liste des utilisateurs
function updateUsersList(users) {
  const usersList = document.querySelector(".users-list");
  if (!usersList) return;

  usersList.innerHTML = "";

  users.forEach((user) => {
    const userEl = document.createElement("div");
    userEl.classList.add("user-item");
    userEl.textContent = user.name;
    userEl.dataset.userId = user.id;

    // ✅ Cliquer sur un utilisateur pour ouvrir le chat
    userEl.addEventListener("click", () => {
      userEl.classList.remove("has-notification");
      handleChatClick(null, user.id, user.name); // ✅ Ouvrir le chat
    });

    usersList.appendChild(userEl);
  });
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
  buildSidebar(); // ✅ Affiche les utilisateurs
  buildMain();
}

export {header, main, sideBar, buildHeader, showApp};
