import { renderCreatePost } from "./new-post.js";

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
    <button>Catégories</button>
    <button>Chat</button>
  </nav>

  <div class="header-right">
    <span class="user-icon"> </span>
  </div>`;

	const postBtn = document.getElementById("new-post-btn");
	postBtn.addEventListener("click", renderCreatePost);
}

function showApp() {
	document.getElementById("auth-container").style.display = "none";
	document.getElementById("app-container").style.display = "block";
	buildHeader();
}

export { header, main, sideBar, buildHeader, showApp };
