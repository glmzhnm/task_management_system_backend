package routes

import (
	"task-manager/controllers"
	"task-manager/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/tasks", controllers.CreateTask)
		protected.GET("/tasks", controllers.GetTasks)
		protected.GET("/tasks/:id", controllers.GetTaskByID)
		protected.PUT("/tasks/:id", controllers.UpdateTask)
		protected.PATCH("/tasks/:id/status", controllers.UpdateTaskStatus)
		protected.DELETE("/tasks/:id", controllers.DeleteTask)
		protected.GET("/tasks/status/:status", controllers.GetTasksByStatus)
		protected.GET("/tasks/priority/:priority", controllers.GetTasksByPriority)
		protected.GET("/tasks/search", controllers.SearchTasks)
		protected.DELETE("/tasks/completed", controllers.DeleteCompletedTasks)
	}
}
