package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"task-manager/database"
	"task-manager/models"
	"task-manager/routes"
)

func main() {
	database.ConnectDB()
	database.DB.AutoMigrate(&models.User{}, &models.Task{})
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}
	database.ConnectDB()
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
