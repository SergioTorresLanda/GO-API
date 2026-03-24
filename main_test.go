package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	// 1. Set Gin to Test Mode so it doesn't print unnecessary logs
	gin.SetMode(gin.TestMode)

	// 2. Setup the router exactly like main.go
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Graphene Engine Online"})
	})

	// 3. Create a fake HTTP request to our endpoint
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// 4. Create a ResponseRecorder to act as our fake browser/client
	w := httptest.NewRecorder()

	// 5. Fire the request into the router
	router.ServeHTTP(w, req)

	// 6. Assert that the response is exactly 200 OK
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, but got %d", w.Code)
	}
}