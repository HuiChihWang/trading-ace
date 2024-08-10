package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"trading-ace/src/config"
	"trading-ace/src/database"
	"trading-ace/src/router"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = database.InitDatabase()
	if err != nil {
		log.Fatal(err)
		return
	}

	if config.GetAppConfig().AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.SetupRouter()

	err = r.Run(":8083")
	if err != nil {
		return
	}
}
