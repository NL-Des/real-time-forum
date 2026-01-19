/* Tutos :
    -https://javascript.developpez.com/actu/83882/L-API-Websockets-ce-que-c-est-et-comment-l-utiliser/
    -https://developer.mozilla.org/fr/docs/Web/API/WebSockets_API/Writing_WebSocket_client_applications
    -https://fr.javascript.info/websocket
    -https://websocket.org/guides/languages/javascript/
    -https://www.educative.io/answers/how-to-use-websockets-in-javascript
     */

// Exemple de websocket simple par Copilot
    // Serveur pour le websocket
// server.js
import { WebSocketServer } from "ws";

const wss = new WebSocketServer({ port: 8080 });

wss.on("connection", (socket) => {
  console.log("Client connecté");

  socket.send("Bienvenue !");

  socket.on("message", (msg) => {
    console.log("Message reçu :", msg.toString());
    socket.send("Reçu : " + msg);
  });

  socket.on("close", () => {
    console.log("Client déconnecté");
  });
});

console.log("Serveur WebSocket sur ws://localhost:8080");

// 