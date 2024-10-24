package daointerfaces

import(

)

type IDeserializeLogicsDAO interface {
	GetDeserializeLogicsByCmdId(cmdId int, requestID string)
	GetDeserializeLogicsBySwVersion(swVersion string, requestID string)
}