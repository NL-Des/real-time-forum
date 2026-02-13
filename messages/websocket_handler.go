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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn *websocket.Conn
	Name string
}

var clients = make(map[int]*Client)
var broadcast = make(chan IncomingMsg)

// ✅ Structure pour les messages entrants — ajout du champ Offset
type IncomingMsg struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	Offset     int    `json:"offset"` // ✅ NOUVEAU
	SenderID   int
	SenderName string
}

type OutgoingMsg struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type HistoryMsg struct {
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type OnlineUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var database *sql.DB

func HandleWebSocket(db *sql.DB) http.HandlerFunc {
	database = db

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

		broadcastOnlineUsers()

		defer func() {
			delete(clients, userID)
			fmt.Println("Client déconnecté:", name, "| Total:", len(clients))
			broadcastOnlineUsers()
		}()

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
				// ✅ Passe l'offset reçu du frontend
				sendHistory(conn, userID, incoming.ReceiverID, incoming.Offset)
			}
		}
	}
}

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

func HandleMessages() {
	for {
		msg := <-broadcast

		now := time.Now()

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

		response := OutgoingMsg{
			Type:      "message",
			Sender:    msg.SenderName,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: now.Format(time.RFC3339),
		}

		jsonResponse, _ := json.Marshal(response)

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

// ✅ Envoyer l'historique avec LIMIT 10 OFFSET
func sendHistory(conn *websocket.Conn, userID int, otherID int, offset int) {
	rows, err := database.Query(`
        SELECT SenderID, Content, CreatedAt
        FROM messages
        WHERE (SenderID = ? AND ReceiverID = ?)
           OR (SenderID = ? AND ReceiverID = ?)
        ORDER BY CreatedAt DESC
        LIMIT 10 OFFSET ?`,
		userID, otherID,
		otherID, userID,
		offset,
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

	// ✅ Inverser pour afficher du plus ancien au plus récent
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// ✅ Envoyer avec hasMore pour savoir s'il reste des messages
	response, _ := json.Marshal(map[string]interface{}{
		"type":     "message_history",
		"messages": messages,
		"offset":   offset,
		"has_more": len(messages) == 10, // ✅ s'il y a 10 résultats, il y en a peut-être plus
	})

	conn.WriteMessage(websocket.TextMessage, response)
}
