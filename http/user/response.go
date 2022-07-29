package user

type UserResponse struct {
	EncID    string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
