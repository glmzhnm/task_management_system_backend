package routes

import (
	"task-manager/controllers"
	"task-manager/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.Use(middlewares.LoggerMiddleware())
	r.Use(middlewares.RateLimiterMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "task-manager"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile)
		protected.PATCH("/profile/password", controllers.ChangePassword)
		protected.GET("/profile/stats", controllers.GetProfileStats)
		protected.DELETE("/profile", controllers.DeleteAccount)

		protected.POST("/tasks", controllers.CreateTask)
		protected.GET("/tasks", controllers.GetTasks)
		protected.GET("/tasks/search", controllers.SearchTasks)
		protected.GET("/tasks/stats", controllers.GetTaskStats)
		protected.GET("/tasks/due", controllers.GetTasksDue)
		protected.GET("/tasks/status/:status", controllers.GetTasksByStatus)
		protected.GET("/tasks/priority/:priority", controllers.GetTasksByPriority)
		protected.GET("/tasks/:id", controllers.GetTaskByID)
		protected.PUT("/tasks/:id", controllers.UpdateTask)
		protected.PATCH("/tasks/:id/status", controllers.UpdateTaskStatus)
		protected.PATCH("/tasks/:id/priority", controllers.UpdateTaskPriority)
		protected.DELETE("/tasks/:id", controllers.DeleteTask)
		protected.DELETE("/tasks/completed", controllers.DeleteCompletedTasks)
	}
}
