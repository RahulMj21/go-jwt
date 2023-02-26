package routes

import (
	"github.com/RahulMj21/go-jwt/controllers"
	"github.com/RahulMj21/go-jwt/middlewares"
	"github.com/gin-gonic/gin"
)

func AllRoutes(router *gin.RouterGroup) {
	// healthcheck route
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Radha Radha..",
		})
	})

	// Auth routes
	router.POST("/signup", controllers.Signup)
	router.POST("/signin", controllers.Signin)
	router.GET("/me", middlewares.CheckAuth, controllers.GetLoggedInUser)
}
