package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/gostartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/check-email", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()

	// userService.RegisterUser(userInput)
}

// func handler(c *gin.Context) {
// 	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	var users []user.User
// 	db.Find(&users)

// 	c.JSON(http.StatusOK, users)
// }
