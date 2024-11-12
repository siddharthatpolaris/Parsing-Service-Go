package serviceinterfaces

import (
	"parsing-service/apps/decoder/models"
)

type IDecoderService interface {
	GetCommandMappingFromCmdID(cmdID int, requestID string) []models.CommandMapping
	GetCommandMappingFromCmdName(cmdName string, requestID string) []models.CommandMapping
	GetCommandMappingFromSwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.CommandMapping
	GetDeserializeLogicsFromCmdId(cmdID int, requestID string) []models.DeserializeLogics
	GetDeserializeLogicsFromSwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.DeserializeLogics
}
