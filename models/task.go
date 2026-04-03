package models

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Status      string `gorm:"default:'pending'" json:"status"`
	Priority    int    `gorm:"default:1" json:"priority"`
}
