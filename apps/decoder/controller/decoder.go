package controller

import (
	services "parsing-service/apps/decoder/service_interfaces"
	"parsing-service/pkg/logger"
)

type decoderController struct {
	decoderService services.IDecoderService
	logger logger.ILogger
}

func NewDecoderController(
	decoderServiceIntf services.IDecoderService,
	logger logger.ILogger,
) *decoderController{
	return &decoderController{
		decoderService: decoderServiceIntf,
		logger: logger,
	}
}