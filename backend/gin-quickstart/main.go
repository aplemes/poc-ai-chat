package main

import (
	"log"
	"os"

	"gin-quickstart/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	router := gin.Default()
	router.Use(corsMiddleware())

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: GROQ_API_KEY is not set — all chat requests will fail")
	}
	chatHandler := handlers.NewChatHandler(apiKey)

	api := router.Group("/api")
	{
		api.POST("/chat/message", chatHandler.SendMessage)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	origin := os.Getenv("ALLOWED_ORIGIN")
	if origin == "" {
		origin = "*"
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Expose-Headers", "X-Session-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
