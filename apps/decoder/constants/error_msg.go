package constants

const (
	ERROR_IN_CREATING_CONSUMER    = "error occurred while creating fetch data Consumer: %v"
	ERR_SUBSCRIBING_TO_TOPIC                       = "error occurred while subscribing to topic: %v || Error: %v"
	ERR_CONSUMING_FROM_KAFKA                       = "error consuming from Kafka: %v"
	ERR_UNMARSHALING_MESSAGE                       = "error unmarshaling message: %v"
	ERR_COMMITTING_OFFSET_SYNC                     = "error occurred while committing offset sync: %v, offset: %v"
	ERROR_IN_GETTING_CMD_ID							= "error occured while getting command id: %v"
	ERROR_IN_GETTING_METER_IP 						= "error occured while getting meter ip: %v"
)