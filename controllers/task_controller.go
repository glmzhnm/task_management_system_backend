package controllers

import (
	"fmt"
	"net/http"
	"task-manager/database"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var notifierClient = resty.New()

func notifyService(message string) {
	resp, err := notifierClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"message": "%s"}`, message)).
		Post("http://notifier:8081/notify")
	if err != nil {
		fmt.Printf("Notifier error: %v\n", err)
	} else {
		fmt.Printf("Notifier response: %s\n", resp.String())
	}
}

func CreateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UserID = userID
	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	go notifyService(fmt.Sprintf("New task created: %s", task.Title))
	c.JSON(http.StatusCreated, task)
}

func GetTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
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
	task.DueDate = updateData.DueDate

	database.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func UpdateTaskStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required,oneof=pending in_progress completed cancelled"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&task).Update("status", statusUpdate.Status)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	result := database.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).Delete(&models.Task{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func GetTasksByStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var tasks []models.Task
	database.DB.Where("user_id = ? AND status = ?", userID, c.Param("status")).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTasksByPriority(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var tasks []models.Task
	database.DB.Where("user_id = ? AND priority = ?", userID, c.Param("priority")).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func SearchTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}
	var tasks []models.Task
	database.DB.Where("user_id = ? AND (title ILIKE ? OR description ILIKE ?)",
		userID, "%"+query+"%", "%"+query+"%").Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func DeleteCompletedTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	result := database.DB.Where("user_id = ? AND status = ?", userID, "completed").Delete(&models.Task{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete completed tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Deleted %d completed tasks", result.RowsAffected),
		"deleted": result.RowsAffected,
	})
}

func GetTaskStats(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	type StatRow struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var stats []StatRow
	database.DB.Model(&models.Task{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&stats)

	var total int64
	database.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"by_status": stats,
	})
}
func GetTasksDue(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var tasks []models.Task

	if err := database.DB.Where("user_id = ? AND due_date IS NOT NULL", userID).Order("due_date asc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
func UpdateTaskPriority(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskID := c.Param("id")

	var input struct {
		Priority int `json:"priority" binding:"required,min=1,max=5"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	database.DB.Model(&task).Update("priority", input.Priority)
	c.JSON(http.StatusOK, gin.H{"message": "Priority updated successfully", "task": task})
}
