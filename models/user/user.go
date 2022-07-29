package user

import (
	"github.com/Artexus/api-matthew-backend/utils/aes"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (u User) EncID() string {
	return aes.EncryptID(int(u.ID))
}

func (u User) TableName() string {
	return "users"
}
