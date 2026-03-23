package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Trade represents a single trading execution in the Graphene protocol
// We added GORM tags to tell the database how to store this.
type Trade struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Amount    float64   `json:"amount"`
	Side      string    `json:"side"` // "buy" or "sell"
	Timestamp time.Time `json:"timestamp"`
}

// Global database variable
var db *gorm.DB

func initDB() {
	// Docker passes this URL from the docker-compose.yml environment variables!
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback just in case we run it outside of Docker
		dsn = "host=localhost user=graphene_admin password=supersecretpassword dbname=graphene_trading port=5432 sslmode=disable"
	}

	var err error
	// Open the connection to Postgres
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// The Magic Trick: AutoMigrate automatically looks at our Trade struct 
	// and creates the exact SQL table for it if it doesn't exist yet.
	db.AutoMigrate(&Trade{})
	log.Println("Database connected and migrated successfully!")
}

func main() {
	// 1. Connect to Postgres
	initDB()

	// 2. Initialize the Gin router
	router := gin.Default()

	// 3. Setup Routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/trades", getTrades)
		v1.POST("/trades", createTrade)
	}

	// 4. Start Server
	router.Run(":8080")
}

func getTrades(c *gin.Context) {
	var trades []Trade
	// Fetch all trades from the Postgres database
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

	// Save the new trade permanently into Postgres
	db.Create(&newTrade)

	// TODO: Go WebSocket broadcast will go here!

	c.IndentedJSON(http.StatusCreated, newTrade)
}