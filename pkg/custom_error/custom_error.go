package customError

import (
	"parsing-service/constants"
)

type CustomError struct {
	errorField     error
	errorCode      constants.CustomErrorCode
	httpStatusCode int
}

func NewCustomError(err error, errorCode constants.CustomErrorCode, httpErrorCode int) *CustomError {
	return &CustomError{errorField: err, errorCode: errorCode, httpStatusCode: httpErrorCode}
}

func (customerror *CustomError) GetErrorField() error {
	return customerror.errorField
}

func (customerror *CustomError) GetErrorCode() string {
	return string(customerror.errorCode)
}

func (customerror *CustomError) GetHttpStatusCode() int {
	return customerror.httpStatusCode
}
