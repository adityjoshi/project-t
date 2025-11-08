package main

import (
	"fmt"
	"log"
	"os"
	"synapse/internal/db"
	"synapse/internal/handlers"
	"synapse/internal/repository"
	"synapse/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize databases
	if err := db.InitPostgres(); err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer db.Pool.Close()

	if err := db.CreateSchema(); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	if err := db.InitChroma(); err != nil {
		log.Printf("Warning: Failed to initialize ChromaDB: %v", err)
		log.Println("ChromaDB might not be running. Start it with: chroma run --path ./chroma_db")
	}

	// Check AI provider
	aiProvider := os.Getenv("AI_PROVIDER")
	if aiProvider == "" || aiProvider == "gemini" {
		geminiKey := os.Getenv("GEMINI_API_KEY")
		if geminiKey == "" {
			log.Println("Warning: GEMINI_API_KEY not set. AI features will not work.")
			log.Println("Get a free API key from: https://makersuite.google.com/app/apikey")
		} else {
			log.Println("Using Google Gemini (free) for AI features")
		}
	}

	// Initialize services
	aiService := services.NewAIService()
	itemRepo := repository.NewItemRepository(db.Pool)
	relationRepo := repository.NewRelationRepository(db.Pool)
	itemService := services.NewItemService(itemRepo, aiService)
	searchService := services.NewSearchService(aiService, itemRepo)
	relationService := services.NewRelationService(itemRepo, relationRepo, aiService)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemService, relationService)
	searchHandler := handlers.NewSearchHandler(searchService)

	// Setup router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Items
		api.POST("/items", itemHandler.CreateItem)
		api.GET("/items", itemHandler.GetAllItems)
		api.GET("/items/:id", itemHandler.GetItem)
		api.DELETE("/items/:id", itemHandler.DeleteItem)
		api.GET("/items/:id/related", itemHandler.GetRelatedItems)

		// Search
		api.GET("/search", searchHandler.Search)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

