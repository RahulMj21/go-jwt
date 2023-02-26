package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/RahulMj21/go-jwt/initializers"
	"github.com/RahulMj21/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(c *gin.Context) {
	accessToken, accessTokenErr := c.Cookie("access_token")
	refreshToken, refreshTokenErr := c.Cookie("refresh_token")
	if accessTokenErr != nil || refreshTokenErr != nil || len(accessToken) == 0 || len(refreshToken) == 0 {
		c.AbortWithStatus(401)
		return
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check if the token expires
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(401)
			return
		}
		// find the user
		user := models.User{}
		initializers.DB.First(&user, claims["user_id"])
		if user.ID == 0 {
			c.AbortWithStatus(401)
			return
		}
		// attach to the req
		c.Set("user", user)
		// continue
		c.Next()
	} else {
		c.AbortWithStatus(401)
		return
	}
}
