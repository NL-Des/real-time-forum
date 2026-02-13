let currentOffset = 0;
let currentReceiverID = null;
let isLoadingHistory = false;
let hasMoreMessages = true;
let currentChatUserId = null;

export function handleChatClick(e) {
  console.log("clique chat");
  e.preventDefault();

  const main = document.querySelector("#main-content");

  main.innerHTML = `<h2>Message</h2>
   <div class="messages">
     <div class="users-list"></div>
   </div>`;

  currentChatUserId = null;

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
        userEl.dataset.userId = user.id;

        userEl.addEventListener("click", () => {
          userEl.classList.remove("has-notification");
          currentChatUserId = user.id;
          openConversation(main, ws, user.id, user.name);
        });

        usersList.appendChild(userEl);
      });
    }

    // ✅ Message reçu en temps réel
    if (data.type === "message") {
      const senderId = data.sender_id;

      if (currentChatUserId === senderId) {
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
      } else {
        const userItem = document.querySelector(
          `.user-item[data-user-id="${senderId}"]`,
        );
        if (userItem) {
          userItem.classList.add("has-notification");
        }
      }
    }

    // ✅ Historique des messages avec pagination
    if (data.type === "message_history") {
      const receivedDiv = document.querySelector(".message-received");
      if (!receivedDiv) return;

      isLoadingHistory = false;
      hasMoreMessages = data.has_more;

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
        receivedDiv.scrollTop = receivedDiv.scrollHeight;
      } else {
        const previousHeight = receivedDiv.scrollHeight;

        if (data.messages) {
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
  currentOffset = 0;
  currentReceiverID = userId;
  currentChatUserId = userId;
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

  // ✅ Scroll en haut → charger les anciens messages
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

// ✅ Ajouter un message EN BAS
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

// ✅ Ajouter un message EN HAUT
function prependMessage(container, sender, content, createdAt, isMine) {
  const msgEl = document.createElement("div");
  msgEl.classList.add("msg-bubble", isMine ? "msg-sent" : "msg-received");
  msgEl.innerHTML = `
    <span class="msg-time">-${new Date(createdAt).toLocaleTimeString()}-</span>
    <strong>${sender}:</strong>
    <p>${content}</p>`;
  container.prepend(msgEl);
}
