package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"log"
	"trading-ace/src/config"
	"trading-ace/src/contract"
	"trading-ace/src/controller"
	"trading-ace/src/database"
	"trading-ace/src/job"
	"trading-ace/src/router"
	"trading-ace/src/scheduler"
)

func main() {
	database.MigrateDB("file://migrations", config.GetAppConfig().Database)

	job.SetUpJobProcessor()
	defer job.ShutDownJobProcessor()

	ethClient, err := ethclient.Dial(config.GetAppConfig().EthereumNode.SocketUrl)
	log.Printf("Connected to Ethereum Node: %s\n", config.GetAppConfig().EthereumNode.SocketUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	contractAddress := "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"
	uniSwapContract, err := contract.NewUniSwapV2Contract(contractAddress, "abi/uniswapv2.abi.json", ethClient)

	uniSwapContract.ListenSwapEvents(controller.GetUniSwapEventControllerInstance().HandleUniSwapV2Event)

	if config.GetAppConfig().AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.SetupRouter()

	sch, _ := scheduler.SetUpScheduler()
	defer scheduler.ShutDowScheduler(sch)

	err = r.Run(":8083")
	if err != nil {
		return
	}
}
