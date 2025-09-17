package main

import (
	"log"

	"github.com/maulanar/gin-kecilin/config"
	"github.com/maulanar/gin-kecilin/database"
	"github.com/maulanar/gin-kecilin/routes"
	"github.com/maulanar/gin-kecilin/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.Init()
	database.Init()

	//set secret key
	utils.SetJWTKey([]byte(config.SECRETKEY))
	routes.SetRouter(r)

	// Start Server
	r.Run(":" + config.PORT)
	log.Println("Server is running on port:" + config.PORT)
}
