package constants

const (
	ERR_RECORD_NOT_FOUND_MSG                      = "record not found"
	ERR_COMMAND_NOT_FOUND                         = "command not found"
	ERR_DEVICE_NOT_FOUND                          = "device not found"
	ERR_COMMAND_NOT_APPLICABLE                    = "command not applicable on given devices"
	ERR_IDENTIFIER_TYPE_NOT_SUPPORTED             = "identifierType Not Supported"
	ERR_IDENTIFIERS_NOT_SUPPORTED                 = "one request per Device is allowed"
	ERR_IDENTIFIERS_NULL_DEVICE_COMMANDS_NOT_NULL = "identifiers cannot be null and deviceCommands should be null"
	ERR_DATE_ARGS_MISSING_FIELDS                  = "both 'from' and 'to' fields are required for date args type"
	ERR_DATE_ARGS_MISSING_FROM_FIELD              = "from field is required for date args type"
	ERR_CATEGORY_NOT_FOUND                        = "category not found"
	ERR_DEVICE_NOT_VALID                          = "these Devices not valid: %v"
	RECORD_NOT_FOUND_MSG                          = "record not found"
	ERR_RECORD_COUNT                              = "error while counting records"
	ERR_STATUS_NOT_FOUND                          = "invalid command status"
	ERR_PROTOCOL_NOT_FOUND                        = "invalid Protocol type"
)
