package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"sync/atomic"

)


type Trade struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}

// 0. The DB
var db *gorm.DB

// 1. Define the pristine, single-threaded, lock-free Hub
type Hub struct {
	// Registered clients.
	clients map[*websocket.Conn]bool
	// Inbound messages from the HTTP POST handler.
	broadcast chan Trade // We create a buffered channel that holds Values not POINTERS to Trades.
	// Register requests from new React Native terminals.
	register chan *websocket.Conn
	// Unregister requests from disconnected terminals.
	unregister chan *websocket.Conn
}

// 2. Instantiate the global matrix Hub (Connected Clients & broadcast channel)
var GrapheneHub = Hub{
	broadcast:  make(chan Trade, 1000), // A buffer of 1000 means it can hold 1,000 pending broadcasts before it blocks.
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
	clients:    make(map[*websocket.Conn]bool),
}

// 3. The Master Goroutine (The Background Worker)
// This function starts exactly once when the server boots.
func (h *Hub) Run() {
	// This infinite loop runs in the background forever
	for {
		select {
		// 1. A new user opens the app
		case client := <-h.register:
			h.clients[client] = true
			log.Println("[Hub] New terminal connected to the matrix.")

		// 2. A user closes the app
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
				log.Println("[Hub] Terminal disconnected.")
			}

		// 3. A new trade pointer drops in from the POST handler
		case tradePointer := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteJSON(tradePointer)
				if err != nil {
					// If writing fails, clean up the dead connection
					client.Close()
					delete(h.clients, client)
				}
			}
		}
	}
}

// --- WEBSOCKET SETUP ---
// 1. Upgrader: Upgrades HTTP to WebSocket and allows cross-origin requests
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict this to your actual app domains!
	},
}

func initDB() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("CRITICAL ERROR: DATABASE_URL must be set in the environment.")
    }

    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

	err = db.AutoMigrate(Trade{}) 
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    log.Println("Database migrations applied successfully!")
	
    // Configure Connection Pooling
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)           // Max idle connections
    sqlDB.SetMaxOpenConns(100)          // Max total connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Prevent stale connections

    log.Println("Database connection pool established successfully!")
}

func main() {
	initDB()
	// --- CONCURRENCY ENGINE ---
	// Boot the lock-free Hub as a Goroutine BEFORE starting the server
	// We boot our dedicated broadcasting engine as a Goroutine.
	// It sits in the background independent of any HTTP requests.
	go GrapheneHub.Run()

	router := gin.Default()

	// --- NEW WEBSOCKET ROUTE ---
	router.GET("/ws", handleConnections)

	// --- KUBERNETES / GITLAB HEALTH CHECK ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Graphene Engine Online"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/trades", getTrades)
		v1.POST("/trades", PostTrade)
		v1.POST("/startmock", startMockDataGenerator)
		v1.POST("/stopmock", stopMockDataGenerator)
	}

	// --- PRODUCTION-READY BINDING ---
	// Listen on 0.0.0.0 for external cloud connectivity (Fly.io/Docker requirement)
	port := "8080"
	log.Printf("Graphene Engine listening on port %s...", port)

	// Using router.Run("0.0.0.0:8080") explicitly satisfies the proxy requirements
	err := router.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

// 4. The HTTP Handler
func PostTrade(c *gin.Context) {
	// We create a local variable in memory
	var newTrade Trade
	// Parse the JSON from the React Native app
	if err := c.ShouldBindJSON(&newTrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTrade.Timestamp = time.Now()

	// ACID Durability: Save to Postgres immediately.
	// If the server crashes 1ms from now, the money is safe.
	if err := db.Create(&newTrade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute trade"})
		return
	}

	// The Pointer: We pass the memory address (&) of the saved trade into the Channel.
	// We do NOT copy the struct. This saves massive amounts of RAM under heavy load.
	GrapheneHub.broadcast <- newTrade
	// Instantly release the HTTP connection back to the user
	c.JSON(http.StatusCreated, newTrade)
}

// Upgrade the initial HTTP request to a persistent WebSocket connection
func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	// Send the connection pointer to the Hub's register channel. No locks.
	GrapheneHub.register <- ws

	// Listen for disconnects
	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			// If the loop breaks, tell the Hub to unregister the connection
			GrapheneHub.unregister <- ws
			break
		}
	}
}

func getTrades(c *gin.Context) {
	var trades []Trade
	query := db
	// 1. Look for the last known ID from the mobile app
	lastIdQuery := c.Query("last_id")

	if lastIdQuery != "" {
		// Convert the string ID to an integer
		lastIdInt, err := strconv.Atoi(lastIdQuery)
		if err == nil {
			// Tell Postgres to strictly fetch newer IDs
			query = query.Where("id > ?", lastIdInt)
		}
	}
	// 2. Fetch the trades in chronological order
	query.Order("id asc").Find(&trades)

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
	//broadcast <- newTrade

	c.IndentedJSON(http.StatusCreated, newTrade)
}

// Global flag to track state across requests
var isGenerating int32
var quitMock chan struct{}

func startMockDataGenerator(c *gin.Context) {

	basePrices := map[string]float64{
		"BTC": 60500.20, "ETH": 1550.50, "SOL": 75.30, "ADA": 0.17,
		"BCH": 450.00, "DOGE": 0.16, "TON": 6.20, "XRP": 0.52,
		"BNB": 410.00, "TRON": 0.11, "SUI": 1.45, "XMR": 345.80,
	}

	if !atomic.CompareAndSwapInt32(&isGenerating, 0, 1) {
		c.JSON(http.StatusConflict, gin.H{"error": "Mock data generator loop is already active"})
		return
	}

	// 1. Re-initialize the safety brakes
	quitMock = make(chan struct{})

	go func() {
		symbols := make([]string, 0, len(basePrices))
		for k := range basePrices {
			symbols = append(symbols, k)
		}
		sides := []string{"buy", "sell"}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for {
			symbol := symbols[r.Intn(len(symbols))]
			symbol2 := symbols[r.Intn(len(symbols))]
			
			volatility := (r.Float64() * 0.01) - 0.005
			
			// 2. Build Trade 1
			newTrade := Trade{
				Symbol:    symbol,
				Price:     basePrices[symbol] * (1.0 + volatility),
				Amount:    0.1 + (r.Float64() * 4.9 * r.Float64() * 100),
				Side:      sides[r.Intn(2)],
				Timestamp: time.Now(),
			}

			// 3. Build Trade 2
			newTrade2 := Trade{
				Symbol:    symbol2,
				Price:     basePrices[symbol2] * (1.0 + volatility),
				Amount:    0.1 + (r.Float64() * 4.9 * r.Float64() * 1000),
				Side:      sides[r.Intn(2)],
				Timestamp: time.Now(),
			}

			// 4. Save BOTH to the database so they get valid, unique IDs
			if err := db.Create(&newTrade).Error; err != nil {
				println("DB Error Trade 1:", err.Error())
			}
			if err := db.Create(&newTrade2).Error; err != nil {
				println("DB Error Trade 2:", err.Error())
			}

			// 5. Broadcast BOTH by value (NO AMPERSANDS)
			GrapheneHub.broadcast <- newTrade 
			GrapheneHub.broadcast <- newTrade2 

			// 6. The Brakes (Replaces time.Sleep)
			select {
			case <-quitMock:
				return 
			case <-time.After(700 * time.Millisecond):
				// Loop continues
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{"status": "High-speed multi-trade mock engine initialized"})
}

func stopMockDataGenerator(c *gin.Context) {
	// If it is running (1), swap it to stopped (0)
	if atomic.CompareAndSwapInt32(&isGenerating, 1, 0) {
		// Closing the channel instantly fires the case <-quitMock in the select block
		close(quitMock) 
		c.JSON(http.StatusOK, gin.H{"status": "Mock data generation safely stopped."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Generator is not currently running."})
	}
}




