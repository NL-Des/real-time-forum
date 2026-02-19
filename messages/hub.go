package messages

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// La connexion HTTP est "Upgradée" en une connexion WebSocket (WS).
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Initialisation de la BDD.
func Init(database *sql.DB) {
	db = database
}

var db *sql.DB

// Chaque nouvelle discussion va générer une nouvelle connexion.
// La fonction sera rutilisée et recréée pour chaque discussion.
func WsHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade de la connexion.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading: %v \n", err)
		return
	}
	defer conn.Close()

	// Récupération du cookie de session de l'envoyeur, plus précisément son numéro de session.
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("No session cookie found.\n")
		conn.Close()
		return
	}

	// On regarde en BDD si l'utilisateur est lié au token.
	var senderID string
	err = db.QueryRow(`
		SELECT u.UserName
		FROM users u
		JOIN session s ON s.UserID = u.id
		WHERE s.Token = ? AND s.ExpiresAt > CURRENT_TIMESTAMP
	`, cookie.Value).Scan(&senderID)
	if err != nil {
		log.Printf("Invalid or expired session %s \n", err)
		conn.Close()
		return
	}

	// On récupère la personne qui doit recevoir le message.
	// Depuis l'url directement (ex : const socket = new WebSocket(`ws://localhost:8080/ws?receiver=miaoutest`);)
	receiverID := r.URL.Query().Get("receiver")
	if receiverID == "" {
		log.Printf("No receiver specified : %s", err)
		conn.Close()
		return
	}

	// Variable créée pour vérifier périodiquement si il y a de nouveaux messages à afficher.
	var lastChecked = time.Now() //.Format("2026-01-02 15:00:00")

	// Goroutine d'écriture des messages.
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				WriteMessagesFromBddToUserScreen(db, conn, senderID, receiverID, &lastChecked)
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Boucle de réception des messages.
	for {
		// Le code s'arrête ici et attend un message d'un client.
		// Réception du message de l'utilisateur en []byte.
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error %s, connexion close by the client or by mistake", err)
			close(done)
			break // Ainsi, on sort de la boucle, ce qui déclenche le "defer conn.Close()"
		}
		if IsMessageTypeValid(mt, conn) != true {
			continue
		}
		if SafetyMessageRateSending2Seconds(db, senderID, conn) != true {
			continue
		}
		if IsMessageNotEmpty(message, conn) != true {
			continue
		}
		if IsMessageTooTall(message, conn) != true {
			continue
		}
		err = CreateMessageInBDD(db, message, senderID, receiverID)
		if err != nil {
			log.Printf("Error %s when writing the new message in the BDD", err)
		}
		continue // Pour pouvoir passer au message suivant.
	}
}
