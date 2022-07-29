package constant

import "errors"

var (
	ErrExistUsername = errors.New("username exist")
	ErrExistEmail    = errors.New("email exist")

	ErrInvalid               = errors.New("invalid")
	ErrInvalidID             = errors.New("invalid id")
	ErrUsernameInvalid       = errors.New("username invalid")
	ErrEmailInvalid          = errors.New("email invalid")
	ErrDescriptionTitleEmpty = errors.New("description and title cannot be empty")
	ErrMaximumDescription    = errors.New("description cannot be more than 255 character")
	ErrPasswordInvalid       = errors.New("password invalid")

	ErrUserNotExist  = errors.New("user is not exist")
	ErrWrongPassword = errors.New("wrong password")

	ErrEncodeFailed        = errors.New("encode failed")
	ErrGenerateTokenFailed = errors.New("generate token failed")

	ErrConvertFailed = errors.New("failed to convert")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenInvalid  = errors.New("token invalid")
)
