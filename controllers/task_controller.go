package controllers

import (
	"fmt"
	"net/http"
	"task-manager/database"
	"task-manager/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	database.DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	var task models.Task
	if err := database.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	if err := database.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var updateData models.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Title = updateData.Title
	task.Description = updateData.Description
	task.Status = updateData.Status
	task.Priority = updateData.Priority

	database.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func UpdateTaskStatus(c *gin.Context) {
	var task models.Task
	if err := database.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&task).Update("status", statusUpdate.Status)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	if err := database.DB.Delete(&models.Task{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func GetTasksByStatus(c *gin.Context) {
	var tasks []models.Task
	database.DB.Where("status = ?", c.Param("status")).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTasksByPriority(c *gin.Context) {
	var tasks []models.Task
	database.DB.Where("priority = ?", c.Param("priority")).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func SearchTasks(c *gin.Context) {
	query := c.Query("q")
	var tasks []models.Task
	database.DB.Where("title ILIKE ?", "%"+query+"%").Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func DeleteCompletedTasks(c *gin.Context) {
	result := database.DB.Where("status = ?", "completed").Delete(&models.Task{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete completed tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted %d completed tasks", result.RowsAffected)})
}
