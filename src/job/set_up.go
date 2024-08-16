package job

import (
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"trading-ace/src/config"
)

var jobClient Client
var server *asynq.Server

func GetClientInstance() Client {
	return jobClient
}

func SetUpJobProcessor() {
	redisConfig := config.GetAppConfig().Redis.Job

	redisClientOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		DB:   redisConfig.Database,
	}

	jobClient = NewClient(redisClientOpt)

	server = asynq.NewServer(&redisClientOpt, asynq.Config{})

	mux := asynq.NewServeMux()
	mux.Handle(string(TypeUniSwapTransaction), NewUniSwapTransactionProcessor())

	go func() {
		if err := server.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()
}

func ShutDownJobProcessor() {
	err := jobClient.Close()
	if err != nil {
		log.Printf("failed to close client: %v", err)
	}

	server.Shutdown()
}
