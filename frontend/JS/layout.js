import {renderCreatePost} from "./new-post.js";
import {Logout} from "./authentication.js";

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
  <div class="profile-section">
      <p id="welcome-message"></p>
      <img src="./frontend/img/profil.gif" alt="image profil" class="profil-icon">
      <button id="logoutBtn">Déconnexion</button>
    </div>
`;

  const postBtn = document.getElementById("new-post-btn");
  postBtn.addEventListener("click", renderCreatePost);
  const logoutBtn = document.getElementById("logoutBtn");
  logoutBtn.addEventListener("click", Logout);
}

function showApp() {
  document.getElementById("auth-container").style.display = "none";
  document.getElementById("app-container").style.display = "contents";
  buildHeader();
}

export {header, main, sideBar, buildHeader, showApp};
