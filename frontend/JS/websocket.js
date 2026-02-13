// Connexion WebSocket partagée globalement
let ws = null;
let reconnectInterval = null;

export function initWebSocket() {
  if (ws && ws.readyState === WebSocket.OPEN) {
    console.log("✅ WebSocket déjà connecté");
    return ws;
  }

  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("✅ Connecté au WebSocket");
    clearInterval(reconnectInterval);
  };

  ws.onclose = function () {
    console.log("❌ Déconnecté du WebSocket");
    // Reconnexion automatique après 3 secondes
    reconnectInterval = setInterval(() => {
      console.log("🔄 Tentative de reconnexion...");
      initWebSocket();
    }, 3000);
  };

  ws.onerror = function (error) {
    console.error("⚠️ Erreur WebSocket:", error);
  };

  return ws;
}

export function getWebSocket() {
  return ws;
}
