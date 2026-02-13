let currentOffset = 0;
let currentReceiverID = null;
let isLoadingHistory = false;
let hasMoreMessages = true;

export function handleChatClick(e) {
  console.log("clique chat");
  e.preventDefault();

  const main = document.querySelector("#main-content");

  main.innerHTML = `<h2>Message</h2>
   <div class="messages">
     <div class="users-list"></div>
   </div>`;

  const ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("Connecté au WebSocket");
  };

  ws.onmessage = function (event) {
    const data = JSON.parse(event.data);

    // ✅ Liste des utilisateurs en ligne
    if (data.type === "online_users") {
      const usersList = document.querySelector(".users-list");
      if (!usersList) return;
      usersList.innerHTML = "";

      data.users.forEach((user) => {
        const userEl = document.createElement("div");
        userEl.classList.add("user-item");
        userEl.textContent = user.name;
        userEl.dataset.id = user.id;
        usersList.appendChild(userEl);

        userEl.addEventListener("click", () => {
          openConversation(main, ws, user.id, user.name);
        });
      });
    }

    // ✅ Message reçu en temps réel
    if (data.type === "message") {
      const receivedDiv = document.querySelector(".message-received");
      if (receivedDiv) {
        appendMessage(
          receivedDiv,
          data.sender,
          data.content,
          data.created_at,
          false,
        );
      }
    }

    // ✅ Historique des messages avec pagination
    if (data.type === "message_history") {
      const receivedDiv = document.querySelector(".message-received");
      if (!receivedDiv) return;

      isLoadingHistory = false;
      hasMoreMessages = data.has_more;

      // ✅ Premier chargement (offset 0) → remplacer le contenu
      if (data.offset === 0) {
        receivedDiv.innerHTML = "";
        if (data.messages) {
          data.messages.forEach((msg) => {
            appendMessage(
              receivedDiv,
              msg.sender,
              msg.content,
              msg.created_at,
              msg.is_mine,
            );
          });
        }
        // ✅ Scroll tout en bas
        receivedDiv.scrollTop = receivedDiv.scrollHeight;
      } else {
        // ✅ Chargement d'anciens messages → ajouter EN HAUT
        const previousHeight = receivedDiv.scrollHeight;

        if (data.messages) {
          // Ajouter les messages au début (plus anciens en haut)
          data.messages.reverse().forEach((msg) => {
            prependMessage(
              receivedDiv,
              msg.sender,
              msg.content,
              msg.created_at,
              msg.is_mine,
            );
          });
        }

        // ✅ Garder la position de scroll
        const newHeight = receivedDiv.scrollHeight;
        receivedDiv.scrollTop = newHeight - previousHeight;
      }
    }
  };

  ws.onclose = function () {
    console.log("❌ Déconnecté du WebSocket");
  };

  ws.onerror = function (error) {
    console.log("⚠️ Erreur WebSocket:", error);
  };
}

// ✅ Ouvrir une conversation
function openConversation(main, ws, userId, userName) {
  // ✅ Reset pagination
  currentOffset = 0;
  currentReceiverID = userId;
  hasMoreMessages = true;
  isLoadingHistory = false;

  main.innerHTML = `<h2>Message</h2>
    <div class="messages">
      <div class="nameOfUser">${userName}</div>
      <div class="message-received"></div>
      <div class="message-content">
        <textarea class="message-sender" placeholder="Écris ton message..."></textarea>
      </div>
      <button class="send-message">Envoyer</button>
    </div>`;

  // ✅ Charger les 10 premiers messages
  ws.send(
    JSON.stringify({
      type: "get_history",
      receiver_id: userId,
      offset: 0,
    }),
  );
  currentOffset = 10;

  // ✅ Détecter scroll en haut → charger les 10 suivants
  const receivedDiv = document.querySelector(".message-received");
  receivedDiv.addEventListener("scroll", () => {
    if (receivedDiv.scrollTop === 0 && !isLoadingHistory && hasMoreMessages) {
      isLoadingHistory = true;

      ws.send(
        JSON.stringify({
          type: "get_history",
          receiver_id: currentReceiverID,
          offset: currentOffset,
        }),
      );

      currentOffset += 10;
    }
  });

  // ✅ Envoyer un message
  document.querySelector(".send-message").addEventListener("click", () => {
    const textarea = document.querySelector(".message-sender");
    const content = textarea.value.trim();
    if (content === "") return;

    ws.send(
      JSON.stringify({
        type: "message",
        receiver_id: userId,
        content: content,
      }),
    );

    const receivedDiv = document.querySelector(".message-received");
    appendMessage(receivedDiv, "Moi", content, new Date().toISOString(), true);

    textarea.value = "";
  });
}

// ✅ Ajouter un message EN BAS (nouveau message)
function appendMessage(container, sender, content, createdAt, isMine) {
  const msgEl = document.createElement("div");
  msgEl.classList.add("msg-bubble", isMine ? "msg-sent" : "msg-received");
  msgEl.innerHTML = `
    <span class="msg-time">-${new Date(createdAt).toLocaleTimeString()}-</span>
    <strong>${sender}:</strong>
    <p>${content}</p>`;
  container.appendChild(msgEl);
  container.scrollTop = container.scrollHeight;
}

// ✅ Ajouter un message EN HAUT (ancien message chargé)
function prependMessage(container, sender, content, createdAt, isMine) {
  const msgEl = document.createElement("div");
  msgEl.classList.add("msg-bubble", isMine ? "msg-sent" : "msg-received");
  msgEl.innerHTML = `
    <span class="msg-time">-${new Date(createdAt).toLocaleTimeString()}-</span>
    <strong>${sender}:</strong>
    <p>${content}</p>`;
  container.prepend(msgEl); // ✅ Ajoute en haut
}
