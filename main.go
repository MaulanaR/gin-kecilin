package main

import (
	"gin/config"
	"gin/db"
	"gin/routes"
	"gin/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.Init()
	db.Init()

	//set secret key
	utils.SetJWTKey([]byte(config.SECRETKEY))
	routes.SetRouter(r)

	// Start Server
	r.Run(":" + config.PORT)
	log.Println("Server is running on port:" + config.PORT)
}
