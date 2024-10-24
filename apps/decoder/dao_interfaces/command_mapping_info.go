package daointerfaces

import(

)

type ICommandMappingDAO interface {
	GetCommandMappingByCmdID(cmdID int, requestID string)
	GetCommandMappingByCmdName(cmdName string, requestID string)
	GetCommandMappingBySwVersion(swVersion string, requestID string)
}
