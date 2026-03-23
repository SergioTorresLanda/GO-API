package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Trade struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}

var db *gorm.DB

// --- WEBSOCKET SETUP ---
// 1. Upgrader: Upgrades HTTP to WebSocket and allows cross-origin requests
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict this to your actual app domains!
	},
}

// 2. The Hub: Keeps track of connected clients and the broadcast channel
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Trade)

// 3. The Broadcaster: Runs in the background and sends trades to all clients
func handleMessages() {
	for {
		// Wait for a new trade to be sent into the broadcast channel
		trade := <-broadcast
		
		// Loop through every connected React Native app
		for client := range clients {
			// Push the JSON trade data down the socket
			err := client.WriteJSON(trade)
			if err != nil {
				log.Printf("Client disconnected: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
// --- END WEBSOCKET SETUP ---

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=graphene_admin password=supersecretpassword dbname=graphene_trading port=5432 sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&Trade{})
	log.Println("Database connected and migrated successfully!")
}

func main() {
	initDB()

	// Start the WebSocket broadcaster in a concurrent Goroutine
	go handleMessages()

	router := gin.Default()

	// --- NEW WEBSOCKET ROUTE ---
	router.GET("/ws", handleConnections)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/trades", getTrades)
		v1.POST("/trades", createTrade)
	}

	router.Run(":8080")
}

// Upgrades the incoming connection and registers the new client
func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	
	// Register the new client connection
	clients[ws] = true
	log.Println("New WebSocket client connected!")
}

func getTrades(c *gin.Context) {
	var trades []Trade
	db.Find(&trades)
	c.IndentedJSON(http.StatusOK, trades)
}

func createTrade(c *gin.Context) {
	var newTrade Trade

	if err := c.BindJSON(&newTrade); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid trade format"})
		return
	}

	newTrade.Timestamp = time.Now()

	// 1. Save to Postgres
	db.Create(&newTrade)

	// 2. Broadcast to all WebSockets!
	// This sends the saved trade into the channel, triggering the handleMessages loop
	broadcast <- newTrade

	c.IndentedJSON(http.StatusCreated, newTrade)
}