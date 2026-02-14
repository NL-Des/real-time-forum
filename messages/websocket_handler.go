package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/auth"
	"sync"
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
	mu   sync.Mutex
}

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.RWMutex
	broadcast    = make(chan IncomingMsg)
)

type IncomingMsg struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	Offset     int    `json:"offset"`
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

		// ✅ METTRE L'UTILISATEUR EN LIGNE DANS LA BASE
		_, err = db.Exec("UPDATE users SET userOnline = 1 WHERE id = ?", userID)
		if err != nil {
			log.Printf("❌ Erreur mise à jour statut connexion: %v\n", err)
		} else {
			log.Printf("✅ %s (ID: %d) est maintenant EN LIGNE\n", name, userID)
		}

		// ✅ Ajouter le client
		clientsMutex.Lock()
		clients[userID] = &Client{Conn: conn, Name: name}
		totalClients := len(clients)
		clientsMutex.Unlock()

		fmt.Printf("✅ Client connecté: %s (ID: %d) | Total: %d\n", name, userID, totalClients)

		// ✅ IMPORTANT : Attendre 100ms que le client soit prêt
		time.Sleep(100 * time.Millisecond)

		// ✅ Diffuser APRÈS le délai
		broadcastOnlineUsers()

		// ✅ Fonction de nettoyage à la déconnexion
		defer func() {
			// ✅ METTRE L'UTILISATEUR HORS LIGNE DANS LA BASE
			_, err := db.Exec("UPDATE users SET userOnline = 0 WHERE id = ?", userID)
			if err != nil {
				log.Printf("❌ Erreur mise à jour statut déconnexion: %v\n", err)
			} else {
				log.Printf("✅ %s (ID: %d) est maintenant HORS LIGNE\n", name, userID)
			}

			clientsMutex.Lock()
			delete(clients, userID)
			totalClients := len(clients)
			clientsMutex.Unlock()

			conn.Close()
			fmt.Printf("❌ Client déconnecté: %s (ID: %d) | Total: %d\n", name, userID, totalClients)

			// ✅ Diffuser APRÈS avoir supprimé le client
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
				sendHistory(conn, userID, incoming.ReceiverID, incoming.Offset)
			}
		}
	}
}

func broadcastOnlineUsers() {
	clientsMutex.RLock()
	users := []OnlineUser{}
	for id, client := range clients {
		users = append(users, OnlineUser{ID: id, Name: client.Name})
	}
	clientsMutex.RUnlock()

	usersJSON, err := json.Marshal(map[string]interface{}{
		"type":  "online_users",
		"users": users,
	})
	if err != nil {
		log.Println("Erreur marshal users:", err)
		return
	}

	fmt.Printf("📡 Diffusion online_users: %d utilisateurs en ligne\n", len(users))
	for _, u := range users {
		fmt.Printf("   - ID: %d, Name: %s\n", u.ID, u.Name)
	}

	// ✅ Envoyer à tous les clients
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()

	for id, client := range clients {
		client.mu.Lock()
		err := client.Conn.WriteMessage(websocket.TextMessage, usersJSON)
		client.mu.Unlock()

		if err != nil {
			log.Printf("❌ Erreur envoi à client %d: %v\n", id, err)
		} else {
			fmt.Printf("✅ Message envoyé au client %d\n", id)
		}
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

		clientsMutex.RLock()
		receiver, ok := clients[msg.ReceiverID]
		clientsMutex.RUnlock()

		if ok {
			receiver.mu.Lock()
			err := receiver.Conn.WriteMessage(websocket.TextMessage, jsonResponse)
			receiver.mu.Unlock()

			if err != nil {
				log.Println("Erreur envoi:", err)
			}
		}
	}
}

func sendHistory(conn *websocket.Conn, userID int, otherID int, offset int) {
	rows, err := database.Query(`
		SELECT SenderID, Content, CreatedAt
		FROM messages
		WHERE (SenderID = ? AND ReceiverID = ?)
		   OR (SenderID = ? AND ReceiverID = ?)
		ORDER BY CreatedAt DESC
		LIMIT 10 OFFSET ?`,
		userID, otherID, otherID, userID, offset,
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
		clientsMutex.RLock()
		if client, ok := clients[senderID]; ok {
			senderName = client.Name
		}
		clientsMutex.RUnlock()

		if senderName == "" {
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

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":     "message_history",
		"messages": messages,
		"offset":   offset,
		"has_more": len(messages) == 10,
	})

	conn.WriteMessage(websocket.TextMessage, response)
}
