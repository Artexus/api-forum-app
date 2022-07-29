package post

import "time"

type GetPostResponse struct {
	ID          int64     `json:"-"`
	UserID      int64     `json:"-"`
	EncID       string    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
