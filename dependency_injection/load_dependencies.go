package dependencyinjection

import (
	decoderServiceIntf "parsing-service/apps/decoder/service_interfaces"
	"parsing-service/pkg/config"
	"parsing-service/pkg/database"
	"parsing-service/pkg/logger"
	"parsing-service/routers"
	"sync"

	"go.uber.org/fx"
)

func LoadDependecies() {
	fx.New(
		fx.Invoke(initializeConnectionsAndConfig),
		configModule,
		loggerModule,
		kafkaFactoy,
		decoderModule,
		routerModule,

		fx.Invoke(StartKafkaConsumer),

		fx.Invoke(serveHttpRequests),
	)
}

func StartKafkaConsumer(
	logger logger.ILogger,
	decoderHandler decoderServiceIntf.IDecoderKafkaConsumerService,
) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		decoderHandler.FetchData()
	}()

}

func initializeConnectionsAndConfig(logger logger.ILogger) {
	if _, err := config.SetupConfig(); err != nil {
		logger.Fatalf("Config Setup error: %s", err)
	}

	err := database.SetUpDbConnection(logger)
	if err != nil {
		logger.Fatalf("erorr while setting DB connection in <initializeConnectionsAndConfig>:%v", err)
	}
}


func serveHttpRequests(r *routers.Routes) {
	logger.GetLogger().Fatalf("%v", r.Router.Run(config.ServerConfig()))
}
