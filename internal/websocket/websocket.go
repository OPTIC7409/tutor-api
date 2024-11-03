package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/OPTIC7409/tutor-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

type Client struct {
	Conn   *websocket.Conn
	UserID uint
}

var (
	clients = make(map[*Client]bool)
	mutex   = &sync.Mutex{}
	DB      *gorm.DB
)

func InitWebSocket(db *gorm.DB) {
	DB = db
}

func New() func(*fiber.Ctx) error {
	return websocket.New(Handler)
}

func Handler(c *websocket.Conn) {
	client := &Client{Conn: c}

	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		c.Close()
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var msg models.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("unmarshal:", err)
			continue
		}

		if err := DB.Create(&msg).Error; err != nil {
			log.Println("create message:", err)
			continue
		}

		broadcastMessage(messageType, message)
	}
}

func broadcastMessage(messageType int, message []byte) {
	mutex.Lock()
	for client := range clients {
		if err := client.Conn.WriteMessage(messageType, message); err != nil {
			log.Println("write:", err)
			client.Conn.Close()
			delete(clients, client)
		}
	}
	mutex.Unlock()
}
