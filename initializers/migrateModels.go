package initializers

import "github.com/RahulMj21/go-jwt/models"

func MigrateModels() {
	DB.AutoMigrate(&models.User{})
}
