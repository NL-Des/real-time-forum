export function updateOnlineUsers(users) {
    const sidebar = document.getElementById("sidebar");
    sidebar.innerHTML = "<h3>Utilisateurs en ligne</h3>";

    const ul = document.createElement("ul");

    users.forEach(u => {
        const li = document.createElement("li");
        li.textContent = u;
        li.style.cursor = "pointer";

        li.onclick = () => {
            import("./chat.js").then(module => {
                module.initChat(u);
            });
        };

        ul.appendChild(li);
    });

    sidebar.appendChild(ul);
}
