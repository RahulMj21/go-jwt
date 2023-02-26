package main

import (
	"github.com/RahulMj21/go-jwt/initializers"
	"github.com/RahulMj21/go-jwt/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDb()
	initializers.MigrateModels()
}

func main() {
	app := gin.New()
	app.Use(gin.Logger())
	router := app.Group("/api")

	routes.AllRoutes(router)

	app.Run()
}
