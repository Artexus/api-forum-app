package constant

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load()

	AccessTokenDuration, _  = strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))
	RefreshTokenDuration, _ = strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION"))

	AccessTokenSignedKey  = os.Getenv("ACCESS_SIGNING_KEY")
	RefreshTokenSignedKey = os.Getenv("REFRESH_SIGNING_KEY")

	AccessTokenInterval  = time.Duration(AccessTokenDuration) * time.Minute
	RefreshTokenInterval = time.Duration(RefreshTokenDuration) * time.Hour

	AESKey          = os.Getenv("AES_ENCRYPT_KEY")
	AESMinLength, _ = strconv.Atoi(os.Getenv("AES_MIN_LENGTH"))
)

const (
	IdentifierTypeUsername = "username"
	IdentifierTypeEmail    = "email"

	UsernameMaxLength        = 12
	PasswordMaxlength        = 15
	MaxPostDescriptionLength = 255

	AccessTokenType  = "access_token"
	RefreshTokenType = "refresh_token"
)
