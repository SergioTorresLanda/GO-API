package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// Trade represents a single trading execution in the Graphene protocol
type Trade struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Side      string    `json:"side"` // "buy" or "sell"
	Timestamp time.Time `json:"timestamp"`
}

// Mock database to hold our initial trades
var trades = []Trade{
	{ID: "t_1001", Symbol: "BTC/USD", Price: 64500.50, Amount: 0.5, Side: "buy", Timestamp: time.Now()},
	{ID: "t_1002", Symbol: "ETH/USD", Price: 3420.10, Amount: 4.2, Side: "sell", Timestamp: time.Now().Add(-1 * time.Minute)},
}

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// API Grouping (Good practice for versioning)
	v1 := router.Group("/api/v1")
	{
		v1.GET("/trades", getTrades)
		v1.POST("/trades", createTrade)
	}

	// Start the server on port 8080
	router.Run(":8080")
}

// getTrades responds with the list of all trades as JSON
func getTrades(c *gin.Context) {
	// c.IndentedJSON formats the JSON nicely for debugging, 
	// use c.JSON in production for speed.
	c.IndentedJSON(http.StatusOK, trades)
}

// createTrade adds a new trade from JSON received in the request body
func createTrade(c *gin.Context) {
	var newTrade Trade

	// Call BindJSON to bind the received JSON to newTrade
	if err := c.BindJSON(&newTrade); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid trade format"})
		return
	}

	// Add a server-side timestamp
	newTrade.Timestamp = time.Now()

	// "Save" to our mock database
	trades = append(trades, newTrade)

	// TODO: Here is where we will trigger the WebSocket broadcast to push 
	// the new trade to the React Native app in real-time.

	c.IndentedJSON(http.StatusCreated, newTrade)
}