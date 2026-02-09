import { main, showApp } from "./layout.js";

let isUsed = false;

function renderCreatePost() {
	console.log("affichage forumalaire post");

	if (isUsed === true) {
		return;
	}
	isUsed = true;

	main.innerHTML = `
    <h2>New post</h2>
    <form id="post-form">
      <input name="title" placeholder="Title" required />
      <textarea name="content" placeholder="Text" required></textarea>
      <div class="select-category">
      <input type="checkbox" name="category" value="1"/>
      <label for="category1"> Category 1 </label>
      <input type="checkbox" name="category" value="2"/>
      <label for="category2"> Category 2 </label>
      <input type="checkbox" name="category" value="3"/>
      <label for="category3"> Category 3 </label>
      <input type="checkbox" name="category" value="4"/>
      <label for="category4"> Category 4 </label>
      </div>
      <button type="submit">Publier</button>
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
	let categoriesId = getSelectedCategories(form);
	console.log(categoriesId);

	const data = {
		title: form.title.value,
		content: form.content.value,
		authorid: 1,
		category_ids: categoriesId,
	};

	const res = await fetch("/post", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(data),
	});

	console.log("data : ", data);

	console.log("fetch fait");

	const result = await res.json();

	if (!res.ok) {
		document.getElementById("error").textContent = result.error;
	} else {
		alert("Post créé !");
		showApp();
	}
}

export { renderCreatePost };
