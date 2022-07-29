/**
 * Created by VoidArtanis on 10/22/2017
 */

package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Artexus/api-matthew-backend/constant"
	"github.com/Artexus/api-matthew-backend/utils/jwt"
	"github.com/gin-gonic/gin"
)

var (
	SigningKey = "$SolidSigningKey$"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		// Check if toke in correct format
		// ie Bearer: xx03xllasx
		b := "Bearer "
		if !strings.Contains(token, b) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Your request is not authorized"})
			c.Abort()
			return
		}
		t := strings.Split(token, b)
		if len(t) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "An authorization token was not supplied"})
			c.Abort()
			return
		}

		// Extract token from access token
		_, err := jwt.ExtractToken(t[1], constant.AccessTokenSignedKey)
		if err != nil {
			log.Println("[ERROR] extract token: ", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
