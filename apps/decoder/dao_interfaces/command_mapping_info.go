package daointerfaces

import "parsing-service/apps/decoder/models"

type ICommandMappingDAO interface {
	GetCommandMappingByCmdID(cmdID int, requestID string) []models.CommandMapping
	GetCommandMappingByCmdName(cmdName string, requestID string)  []models.CommandMapping
	GetCommandMappingBySwVersionAndCmdID(swVersion string, cmdID int, requestID string) []models.CommandMapping
}
