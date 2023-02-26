package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/RahulMj21/go-jwt/initializers"
	"github.com/RahulMj21/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func Signup(c *gin.Context) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	validationErr := Validate.Struct(user)
	if validationErr != nil {
		c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	user.Password = string(hash)

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	accessToken, accessTokenerr := SignToken(user.ID, time.Now().Add(time.Hour).Unix(), os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshToken, refreshTokenerr := SignToken(user.ID, time.Now().Add(time.Hour*24*30).Unix(), os.Getenv("REFRESH_TOKEN_SECRET"))
	if accessTokenerr != nil || refreshTokenerr != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "failed to create tokens"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, 1000*60*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 1000*60*60*24*30, "/", "", false, true)

	c.JSON(201, gin.H{"status": "success", "data": user})
}

func Signin(c *gin.Context) {
	var loginBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=25"`
	}

	if err := c.BindJSON(&loginBody); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	validationErr := Validate.Struct(loginBody)
	if validationErr != nil {
		c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
		return
	}

	user := models.User{}

	initializers.DB.Where("email = ?", loginBody.Email).First(&user)
	if user.ID == 0 {
		c.JSON(404, gin.H{"status": "fail", "message": "wrong email"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password))
	if err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": "wrong password"})
		return
	}

	accessToken, accessTokenerr := SignToken(user.ID, time.Now().Add(time.Hour).Unix(), os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshToken, refreshTokenerr := SignToken(user.ID, time.Now().Add(time.Hour*24*30).Unix(), os.Getenv("REFRESH_TOKEN_SECRET"))
	if accessTokenerr != nil || refreshTokenerr != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "failed to create tokens"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, 1000*60*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 1000*60*60*24*30, "/", "", false, true)

	c.JSON(200, gin.H{"status": "success", "data": user})
}

func GetLoggedInUser(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(200, gin.H{"status": "success", "data": user})
}

func SignToken(userId uint, exp int64, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     exp,
	})

	return token.SignedString([]byte(secret))
}
