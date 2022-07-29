package validate

import (
	"net/mail"
	"strings"
	"unicode"

	"github.com/Artexus/api-matthew-backend/constant"
)

func ValidateRegister(username, email, password string) (err error) {
	err = ValidateUsername(username)
	if err != nil {
		return
	}

	err = ValidateEmail(email)
	if err != nil {
		return
	}

	err = ValidatePassword(password)
	return
}

func ValidateUsername(username string) error {
	if len(username) > constant.UsernameMaxLength ||
		strings.Contains(username, " ") {
		return constant.ErrUsernameInvalid
	}
	return nil
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return constant.ErrEmailInvalid
	}
	return nil
}

func ValidatePassword(password string) error {
	isDigit := false
	isUpper := false
	for _, s := range password {
		if unicode.IsDigit(s) {
			isDigit = true
		} else if unicode.IsUpper(s) {
			isUpper = true
		}
	}

	if len(password) > constant.PasswordMaxlength ||
		!isDigit || !isUpper {
		return constant.ErrPasswordInvalid
	}
	return nil
}
