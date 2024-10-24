package dependencyinjection

import (
	decoderServiceIntf "parsing-service/apps/decoder/service_interfaces"
	"parsing-service/pkg/config"
	"parsing-service/pkg/database"
	"parsing-service/pkg/logger"
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

		fx.Invoke(StartKafkaConsumer),
	)
}

func StartKafkaConsumer(
	logger logger.ILogger,
	decoderHandler decoderServiceIntf.IDecoderService,
) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		decoderHandler.FetchData()
	}()

	wg.Wait()
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
