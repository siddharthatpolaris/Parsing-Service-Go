package dependencyinjection

import (
	"parsing-service/pkg/config"
	"parsing-service/pkg/logger"

	kafkaImpl "parsing-service/apps/kafka/service_impl"
	kafkaIntf "parsing-service/apps/kafka/service_interfaces"

	decoderController "parsing-service/apps/decoder/controller"
	decoderServiceInt "parsing-service/apps/decoder/service_interfaces"
	decoderServices "parsing-service/apps/decoder/services"

	routers "parsing-service/routers"

	"go.uber.org/fx"
)

var configModule = fx.Options(
	fx.Provide(
		config.SetupConfig,
	),
)

var loggerModule = fx.Provide(
	fx.Annotate(
		logger.NewLogger,
		fx.As(new(logger.ILogger)),
	),
)

var kafkaFactoy = fx.Options(
	fx.Provide(
		fx.Annotate(
			kafkaImpl.NewKafkaConsumerFactory,
			fx.As(new(kafkaIntf.IKafkaConsumerFactory)),
		),
	),
)

var decoderModule = fx.Options(
	fx.Provide(
		decoderController.NewDecoderController,
		fx.Annotate(
			decoderServices.NewKafkaConsumerHandler,
			fx.As(new(decoderServiceInt.IDecoderKafkaConsumerService)),
		),
		fx.Annotate(
			decoderServices.NewDecoder,
			fx.As(new(decoderServiceInt.IDecoderService)),
		),
	),
)


var routerModule = fx.Options(
	fx.Provide(
		routers.NewHandler,
	),
)
