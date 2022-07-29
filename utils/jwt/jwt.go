package jwt

import (
	"log"
	"strings"
	"time"

	"github.com/Artexus/api-matthew-backend/constant"
	"github.com/Artexus/api-matthew-backend/utils/aes"
	"github.com/form3tech-oss/jwt-go"
)

type TokenResponse struct {
	EncID     string
	UserID    int
	Username  string
	Email     string
	ExpiredIn int64
}

func (tr TokenResponse) IsExpired() bool {
	return !time.Unix(tr.ExpiredIn, 0).After(time.Now())
}

func GenerateToken(userID, email, username string, expired int64, signedKey string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)
	claims["user_id"] = userID
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = expired
	token.Claims = claims

	tokenString, err = token.SignedString([]byte(signedKey))
	return
}

func ExtractIDToken(token string, signedKey string) (id int, err error) {
	claims := jwt.MapClaims{}

	if strings.Contains(token, "Bearer ") {
		token = strings.Split(token, "Bearer ")[1]
	}
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	if err != nil {
		log.Println(err)
		err = constant.ErrTokenInvalid
		return
	}

	log.Println("ENC USER : ", claims["user_id"].(string))
	id = aes.DecryptID(claims["user_id"].(string))

	return
}

func ExtractToken(token string, signedKey string) (resp TokenResponse, err error) {
	claims := jwt.MapClaims{}

	if strings.Contains(token, "Bearer ") {
		token = strings.Split(token, "Bearer ")[1]
	}
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	if err != nil && strings.Contains(err.Error(), "expired") {
		err = constant.ErrTokenExpired
		return
	} else if err != nil {
		err = constant.ErrTokenInvalid
		return
	}

	resp.EncID = claims["user_id"].(string)
	resp.UserID = aes.DecryptID(resp.EncID)
	resp.Username = claims["username"].(string)
	resp.Email = claims["email"].(string)
	resp.ExpiredIn = int64(claims["exp"].(float64))
	return
}
