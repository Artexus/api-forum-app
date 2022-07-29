package token

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	UserID    int    `gorm:"column:user_id"`
	Type      string `gorm:"column:type"`
	Token     string `gorm:"column:token"`
	ExpiredAt int64  `gorm:"column:expired_at"`
}

func (t Token) TableName() string {
	return "tokens"
}
