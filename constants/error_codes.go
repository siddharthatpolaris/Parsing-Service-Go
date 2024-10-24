package constants

type CustomErrorCode string

const(
	INVALID_REQUEST_BODY        CustomErrorCode = "RS-402"
	INTERNAL_SERVER_ERROR_CODE  CustomErrorCode = "RS-500"
	GENERIC_ERROR_CODE          CustomErrorCode = "RS-501"
	INVALID_JSON_ERROR_CODE     CustomErrorCode = "RS-401"
	RECORD_NOT_FOUND_ERROR_CODE CustomErrorCode = "RS-404"
	INVALID_REQ_PARAMS_CODE     CustomErrorCode = "RS-403"
	BAD_REQUEST_ERROR_CODE      CustomErrorCode = "RS-400"
)