package post

import "github.com/jinzhu/gorm"

type Post struct {
	gorm.Model
	UserID      int    `gorm:"column:user_id"`
	Description string `gorm:"column:description"`
	Title       string `gorm:"column:title"`
	FileID      *int   `gorm:"column:file_id"`
}

func (p Post) TableName() string {
	return "posts"
}
