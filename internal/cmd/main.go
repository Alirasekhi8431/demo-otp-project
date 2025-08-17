package main

import (
	"alirasekhi8431/demo-otp-project/internal/api"
	db "alirasekhi8431/demo-otp-project/internal/db"
	
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectToDb("myuser", "mypassword", "5432")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		gin.Logger()(c)
	})
	db.InsertUser("john_doe",  "some_pass")
	api.SetupRoutes(router)
	router.Run(":8080")

}

