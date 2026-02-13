package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/auth"
	"time"

	"github.com/gorilla/websocket"
)

// Configuration du WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Structure pour un client connecté
type Client struct {
	Conn *websocket.Conn
	Name string
}

// Map pour gérer les clients connectés : userID → Client
var clients = make(map[int]*Client)

// Channel pour transmettre les messages
var broadcast = make(chan IncomingMsg)

// ✅ Structure pour les messages entrants (depuis le frontend)
type IncomingMsg struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	SenderID   int    // rempli côté serveur
	SenderName string // rempli côté serveur
}

// ✅ Structure pour les messages sortants (vers le frontend)
type OutgoingMsg struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// ✅ Structure pour l'historique
type HistoryMsg struct {
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

// Structure pour la liste des utilisateurs en ligne
type OnlineUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Variable globale pour la BDD
var database *sql.DB

// Handler WebSocket
func HandleWebSocket(db *sql.DB) http.HandlerFunc {
	database = db // ✅ Stocker la référence à la BDD

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(int)

		var name string
		err := db.QueryRow("SELECT UserName FROM users WHERE id=?", userID).Scan(&name)
		if err != nil {
			log.Println("Erreur récupération nom:", err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Erreur WebSocket upgrade:", err)
			return
		}
		defer conn.Close()

		clients[userID] = &Client{Conn: conn, Name: name}
		fmt.Println("Client connecté !", name, "| Total:", len(clients))

		// ✅ Envoyer la liste des users en ligne à TOUS
		broadcastOnlineUsers()

		defer func() {
			delete(clients, userID)
			fmt.Println("Client déconnecté:", name, "| Total:", len(clients))
			broadcastOnlineUsers()
		}()

		// Écouter les messages
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("Erreur lecture:", err)
				break
			}

			var incoming IncomingMsg
			if err := json.Unmarshal(msgBytes, &incoming); err != nil {
				log.Println("Erreur JSON:", err)
				continue
			}

			incoming.SenderID = userID
			incoming.SenderName = name

			switch incoming.Type {
			case "message":
				broadcast <- incoming

			case "get_history":
				sendHistory(conn, userID, incoming.ReceiverID)
			}
		}
	}
}

// ✅ Envoyer la liste des utilisateurs en ligne à tous
func broadcastOnlineUsers() {
	users := []OnlineUser{}
	for id, client := range clients {
		users = append(users, OnlineUser{ID: id, Name: client.Name})
	}

	usersJSON, err := json.Marshal(map[string]interface{}{
		"type":  "online_users",
		"users": users,
	})
	if err != nil {
		return
	}

	for _, client := range clients {
		client.Conn.WriteMessage(websocket.TextMessage, usersJSON)
	}
}

// ✅ Sauvegarder en BDD + envoyer au destinataire
func HandleMessages() {
	for {
		msg := <-broadcast

		now := time.Now()

		// ✅ Sauvegarder en BDD
		_, err := database.Exec(`
            INSERT INTO messages (SenderID, ReceiverID, Content, CreatedAt)
            VALUES (?, ?, ?, ?)`,
			msg.SenderID, msg.ReceiverID, msg.Content, now,
		)
		if err != nil {
			log.Println("Erreur insertion message:", err)
			continue
		}

		fmt.Println("Message sauvegardé:", msg.SenderName, "→", msg.ReceiverID, ":", msg.Content)

		// ✅ Construire la réponse
		response := OutgoingMsg{
			Type:      "message",
			Sender:    msg.SenderName,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: now.Format(time.RFC3339),
		}

		jsonResponse, _ := json.Marshal(response)

		// ✅ Envoyer au destinataire s'il est en ligne
		if receiver, ok := clients[msg.ReceiverID]; ok {
			err := receiver.Conn.WriteMessage(websocket.TextMessage, jsonResponse)
			if err != nil {
				log.Println("Erreur envoi:", err)
				receiver.Conn.Close()
				delete(clients, msg.ReceiverID)
			}
		}
	}
}

// ✅ Envoyer l'historique des messages entre 2 utilisateurs
func sendHistory(conn *websocket.Conn, userID int, otherID int) {
	rows, err := database.Query(`
        SELECT SenderID, Content, CreatedAt
        FROM messages
        WHERE (SenderID = ? AND ReceiverID = ?)
           OR (SenderID = ? AND ReceiverID = ?)
        ORDER BY CreatedAt ASC`,
		userID, otherID,
		otherID, userID,
	)
	if err != nil {
		log.Println("Erreur historique:", err)
		return
	}
	defer rows.Close()

	var messages []HistoryMsg

	for rows.Next() {
		var senderID int
		var content, createdAt string
		rows.Scan(&senderID, &content, &createdAt)

		// ✅ Récupérer le nom du sender
		senderName := ""
		if client, ok := clients[senderID]; ok {
			senderName = client.Name
		} else {
			database.QueryRow("SELECT UserName FROM users WHERE id=?", senderID).Scan(&senderName)
		}

		messages = append(messages, HistoryMsg{
			Sender:    senderName,
			SenderID:  senderID,
			Content:   content,
			CreatedAt: createdAt,
			IsMine:    senderID == userID,
		})
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":     "message_history",
		"messages": messages,
	})

	conn.WriteMessage(websocket.TextMessage, response)
}
