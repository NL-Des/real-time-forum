import {main, showApp} from "./layout.js";

let isUsed = false;

function renderCreatePost() {
  console.log("affichage forumalaire post");

  if (isUsed === true) {
    return;
  }
  isUsed = true;

  main.innerHTML = `
    <h2>Nouveau post</h2>
    <form id="post-form">
    <div class="title">
    <label for="title">Titre</label>
      <textarea id="title" name="title" required></textarea>
      </div>
      <div class="message">
        <label for="message">Message</label>
      <textarea id="message" name="message" placeholder="Text" required></textarea>
      </div>
      <div class="select-category">
       <label for="category">Catégorie</label>
      <select id="category" name="category" required>
        <option value="" disabled selected>Choisir une catégorie</option>
        <option value="1">Category 1</option>
        <option value="2">Category 2</option>
        <option value="3">Category 3</option>
        <option value="4">Category 4</option>
      </select>
      </div>
      <div class="newPost">
      <button id="newPost" type="submit">Publier</button>
        </div>
      <p id="error"></p>
    </form>
  `;

  document
    .getElementById("post-form")
    .addEventListener("submit", handleCreatePost);
}

function getSelectedCategories(form) {
  const checked = form.querySelectorAll('input[name="category"]:checked');
  return Array.from(checked).map((cb) => Number(cb.value));
}

async function handleCreatePost(e) {
  isUsed = false;
  e.preventDefault();
  console.log("handleCreatePost");

  const form = e.target;

  const data = {
    title: form.title.value,
    content: form.content.value,
    authorid: 1,
    category_ids: getSelectedCategories(form),
  };
  const res = await fetch("/post", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify(data),
  });

  console.log("fetch fait");

  const result = await res.json();

  if (!res.ok) {
    document.getElementById("error").textContent = result.error;
  } else {
    alert("Post créé !");
    showApp();
  }
}

export {renderCreatePost};
