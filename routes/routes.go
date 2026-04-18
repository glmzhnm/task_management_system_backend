package routes

import (
	"github.com/gin-gonic/gin"
	"task-manager/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/tasks", controllers.CreateTask)
	r.GET("/tasks", controllers.GetTasks)
	r.GET("/tasks/:id", controllers.GetTaskByID)
	r.PUT("/tasks/:id", controllers.UpdateTask)
	r.PATCH("/tasks/:id/status", controllers.UpdateTaskStatus)
	r.DELETE("/tasks/:id", controllers.DeleteTask)
	r.GET("/tasks/status/:status", controllers.GetTasksByStatus)
	r.GET("/tasks/priority/:priority", controllers.GetTasksByPriority)
	r.GET("/tasks/search", controllers.SearchTasks)
	r.DELETE("/tasks/completed", controllers.DeleteCompletedTasks)
}
