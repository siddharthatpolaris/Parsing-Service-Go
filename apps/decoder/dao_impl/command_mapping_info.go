package daoimpl

import (
	"errors"
	"net/http"
	"parsing-service/apps/decoder/models"
	"parsing-service/constants"
	customErrorPkg "parsing-service/pkg/custom_error"
	"parsing-service/pkg/database"
	"parsing-service/pkg/logger"
)

type CommandMappingImpl struct {
	logger logger.ILogger
}

func (commandMappingDao *CommandMappingImpl) GetCommandMappingByCmdID(cmdID int, requestID string) []models.CommandMapping {
	var commandMappingData []models.CommandMapping

	query := database.DB.Model(&models.CommandMapping{}).Where("command_mapping.cmd_id = ?", cmdID)

	err := query.Find(&commandMappingData).Error
	if err != nil {
		commandMappingDao.logger.Errorf("<GetCommandMappingByCmdID> RequestID %v, Error %v", requestID, err)
		customError := customErrorPkg.NewCustomError(
			errors.New(constants.PROCESSING_ERROR),
			constants.INTERNAL_SERVER_ERROR_CODE,
			http.StatusInternalServerError,
		)

		panic(customError)
	}

	return commandMappingData
}


func (commandMappingDao *CommandMappingImpl) GetCommandMappingByCmdName(cmdName string, requestID string) []models.CommandMapping {
	var commandMappingData []models.CommandMapping

	query := database.DB.Model(&models.CommandMapping{}).Where("command_mapping.cmd_name = ?", cmdName)

	err := query.Find(&commandMappingData).Error
	if err != nil {
		commandMappingDao.logger.Errorf("<GetCommandMappingByCmdName> RequestID %v, Error %v", requestID, err)
		customError := customErrorPkg.NewCustomError(
			errors.New(constants.PROCESSING_ERROR),
			constants.INTERNAL_SERVER_ERROR_CODE,
			http.StatusInternalServerError,
		)

		panic(customError)
	}

	return commandMappingData
}



func (commandMappingDao *CommandMappingImpl) GetCommandMappingBySwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.CommandMapping {
	var commandMappingData []models.CommandMapping

	query := database.DB.Model(&models.CommandMapping{}).Where("command_mapping.cmd_id = ?", cmdID)

	query = query.Joins("CommandMappingSwVersion")
	query = query.Joins("SwVersion").Where("sw_version.version = ?", swVersion)

	err := query.Find(&commandMappingData).Error
	if err != nil {
		commandMappingDao.logger.Errorf("<GetCommandMappingBySwVersionAndCmdID> RequestID %v, Error %v", requestID, err)
		customError := customErrorPkg.NewCustomError(
			errors.New(constants.PROCESSING_ERROR),
			constants.INTERNAL_SERVER_ERROR_CODE,
			http.StatusInternalServerError,
		)

		panic(customError)
	}

	return commandMappingData
}


