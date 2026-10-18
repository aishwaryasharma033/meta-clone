package main

import (
	"log"
	"net/http"

	"instagram-backend/config"
	"instagram-backend/controllers"

	"github.com/gin-gonic/gin"
)
func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("PostgreSQL connected successfully!")


	router := gin.Default()

	authController := controllers.AuthController{DB: db}

    router.POST("/signup", authController.Signup)

	router.POST("/login", authController.Login)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Instagram backend is running!",
		})
	})

	router.Run(":8080")
}