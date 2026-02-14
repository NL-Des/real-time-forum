import {renderCreatePost} from "./new-post.js";
import {Logout} from "./authentication.js";
import {handleChatClick} from "./chat.js";
import {initWebSocket, addMessageHandler} from "./websocket.js";

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
  sideBar.innerHTML = `<h2>Utilisateurs</h2>
  <div class="users-list"></div>`;

  // ✅ Initialiser WebSocket AVANT de charger les utilisateurs
  const ws = initWebSocket();

  // ✅ Attendre que le WebSocket soit connecté
  ws.addEventListener("open", () => {
    console.log("✅ WebSocket prêt, chargement utilisateurs...");
    loadAllUsers();
  });

  // ✅ Si déjà connecté, charger immédiatement
  if (ws.readyState === WebSocket.OPEN) {
    loadAllUsers();
  }

  // ✅ Enregistrer le gestionnaire de messages pour la sidebar
  addMessageHandler(handleSidebarMessages);
}

// ✅ Gestionnaire de messages WebSocket pour la sidebar
function handleSidebarMessages(data) {
  // ✅ Mise à jour des utilisateurs en ligne
  if (data.type === "online_users") {
    console.log("📡 [SIDEBAR] Mise à jour utilisateurs:", data.users);
    updateUsersList(data.users);
  }

  // ✅ Notification pour nouveau message
  if (data.type === "message") {
    const userItem = document.querySelector(
      `.user-item[data-user-id="${data.sender_id}"]`,
    );
    if (userItem) {
      userItem.classList.add("has-notification");
    }
  }
}

// ✅ Charger tous les utilisateurs depuis l'API
async function loadAllUsers() {
  const usersList = document.querySelector(".users-list");

  if (!usersList) {
    console.error("❌ .users-list introuvable dans le DOM !");
    return;
  }

  console.log("🔄 Chargement des utilisateurs...");

  try {
    const response = await fetch("/api/users");
    if (!response.ok) throw new Error("Erreur récupération utilisateurs");

    const allUsers = await response.json();
    console.log("✅ Utilisateurs chargés:", allUsers);

    usersList.innerHTML = "";

    if (allUsers.length === 0) {
      console.warn("⚠️ Aucun utilisateur trouvé dans la base");
      usersList.innerHTML = "<p>Aucun utilisateur</p>";
      return;
    }

    allUsers.forEach((user) => {
      const userEl = document.createElement("div");
      userEl.classList.add("user-item");

      // ✅ Par défaut, tous sont hors ligne (classe .offline)
      userEl.classList.add("offline");

      userEl.textContent = user.nickname;
      userEl.dataset.userId = user.id;

      userEl.addEventListener("click", () => {
        userEl.classList.remove("has-notification");
        handleChatClick(null, user.id, user.nickname);
      });

      usersList.appendChild(userEl);
    });

    console.log(
      `✅ ${allUsers.length} utilisateurs affichés (hors ligne par défaut)`,
    );
  } catch (error) {
    console.error("❌ Erreur chargement utilisateurs:", error);
  }
}

// ✅ Mettre à jour les statuts (en ligne/hors ligne)
function updateUsersList(onlineUsers) {
  const usersList = document.querySelector(".users-list");
  if (!usersList) {
    console.error("❌ .users-list introuvable pour mise à jour");
    return;
  }

  // ✅ Créer un Set des IDs en ligne
  const onlineUserIds = new Set(onlineUsers.map((u) => u.id));
  console.log("🟢 Utilisateurs en ligne:", Array.from(onlineUserIds));

  // ✅ Parcourir tous les .user-item et mettre à jour leur statut
  const userItems = usersList.querySelectorAll(".user-item");

  if (userItems.length === 0) {
    console.warn("⚠️ Aucun .user-item trouvé pour mise à jour");
    return;
  }

  userItems.forEach((userEl) => {
    const userId = parseInt(userEl.dataset.userId);

    if (onlineUserIds.has(userId)) {
      // 🟢 En ligne → retirer .offline
      userEl.classList.remove("offline");
      console.log(`🟢 ${userEl.textContent} est EN LIGNE`);
    } else {
      // 🔴 Hors ligne → ajouter .offline
      userEl.classList.add("offline");
      console.log(`🔴 ${userEl.textContent} est HORS LIGNE`);
    }
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
  buildSidebar();
  buildMain();
}

export {header, main, sideBar, buildHeader, showApp};
