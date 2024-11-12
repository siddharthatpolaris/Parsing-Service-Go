package daointerfaces

import "parsing-service/apps/decoder/models"

type IDeserializeLogicsDAO interface {
	GetDeserializeLogicsByCmdId(cmdId int, requestID string) []models.DeserializeLogics
	GetDeserializeLogicsBySwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.DeserializeLogics
}