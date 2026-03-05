// Import de la fonction d’envoi de message
import { sendMessage } from "./sendMessage.js";
import { chatState } from "./chatState.js";

// openConversation Ouvrir une conversation avec un utilisateur
export function openConversation(main, ws, userId, userName) {
	// Log informatif
	console.log(`Ouverture conversation: ${userName} (ID: ${userId})`);

	chatState.currentChatUser = userId;
	chatState.isLoadingHistory = false;
	chatState.hasMoreMessages = true;
	chatState.offset = 0;

	// Construction dynamique de l’interface HTML
	main.innerHTML = `
    <h2>Message avec ${userName}</h2>
    <div class="messages">
      <div class="users-list"></div>
      <div class="conversation">
        <div class="message-received"></div>
        <div class="message-content">
          <textarea class="message-sender" placeholder="Écris ton message..."></textarea>
        </div>
        <button class="send-message">Envoyer</button>
      </div>
    </div>`;

	// Chargement initial de l’historique
	ws.send(
		JSON.stringify({
			type: "get_history", // Type de requête
			receiver_id: userId, // ID du destinataire
			offset: 0, // Premier chargement
		}),
	);

	// On prépare l’offset suivant
	chatState.offset = 10;

	// Gestion du scroll (pagination)
	const receivedDiv = document.querySelector(".message-received");

	receivedDiv.addEventListener("scroll", () => {
		// Si on est en haut
		// ET qu’on ne charge pas déjà
		// ET qu’il reste des messages
		if (
			receivedDiv.scrollTop === 0 &&
			!chatState.isLoadingHistory &&
			chatState.hasMoreMessages
		) {
			console.log("Chargement messages plus anciens...");

			// On bloque les requêtes multiples
			chatState.isLoadingHistory = true;

			// Demande des anciens messages
			ws.send(
				JSON.stringify({
					type: "get_history",
					receiver_id: chatState.currentChatUser,
					offset: chatState.offset,
				}),
			);

			// On prépare l’offset suivant
			chatState.offset += 10;
		}
	});

	// Envoi du message via bouton
	document.querySelector(".send-message").addEventListener("click", () => {
		sendMessage(ws, chatState.currentChatUser);
	});

	// Envoi du message via touche Entrée
	document
		.querySelector(".message-sender")
		.addEventListener("keypress", (e) => {
			// Si touche Entrée sans Shift
			if (e.key === "Enter" && !e.shiftKey) {
				e.preventDefault(); // Empêche retour à la ligne
				sendMessage(ws, chatState.currentChatUser);
			}
		});
}
