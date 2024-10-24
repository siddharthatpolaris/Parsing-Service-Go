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

type DeserializeLogicsImpl struct {
	logger logger.ILogger
}

func (DeserializeLogicsDao *DeserializeLogicsImpl) GetDeserializeLogicsByCmdId(cmdID int, requestID string) []models.DeserializeLogics {
	var deserializeLogicsData []models.DeserializeLogics

	query := database.DB.Model(&models.DeserializeLogics{}).Where("deserialize_logics.cmd_id = ?", cmdID)

	err := query.Find(&deserializeLogicsData).Error
	if err != nil {
		DeserializeLogicsDao.logger.Errorf("<GetDeserializeLogicsByCmdId> RequestID %v, Error %v", requestID, err)
		customError := customErrorPkg.NewCustomError(
			errors.New(constants.PROCESSING_ERROR),
			constants.INTERNAL_SERVER_ERROR_CODE,
			http.StatusInternalServerError,
		)

		panic(customError)
	}

	return deserializeLogicsData
}



func (DeserializeLogicsDao *DeserializeLogicsImpl) GetDeserializeLogicsBySwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.DeserializeLogics {
	var deserializeLogicsData []models.DeserializeLogics

	query := database.DB.Model(&models.DeserializeLogics{}).Where("deserialize_logics.cmd_id = ?", cmdID)

	query = query.Joins("DeserializeLogicSwVersion")
	query = query.Joins("SwVersion").Where("sw_version.version = ?", swVersion)

	err := query.Find(&deserializeLogicsData).Error
	if err != nil {
		DeserializeLogicsDao.logger.Errorf("<GetDeserializeLogicsBySwVersionAndCmdID> RequestID %v, Error %v", requestID, err)
		customError := customErrorPkg.NewCustomError(
			errors.New(constants.PROCESSING_ERROR),
			constants.INTERNAL_SERVER_ERROR_CODE,
			http.StatusInternalServerError,
		)

		panic(customError)
	}

	return deserializeLogicsData
}
