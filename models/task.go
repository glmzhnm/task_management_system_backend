package models

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title" binding:"required"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:'pending'" json:"status"`
	Priority    int       `gorm:"default:1" json:"priority"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
