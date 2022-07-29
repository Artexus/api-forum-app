package token

type InsertTokenRequest struct {
	Type  string `json:"type" validate:"required"`
	Token string `json:"token" validate:"required"`

	UserID    int   `json:"user_id"`
	ExpiredAt int64 `json:"expired_at"`
}

type GenerateTokenRequest struct {
	Token       string `json:"token"`
	RequestType string `json:"request_type"`

	TokenType string `json:"-"`
}
