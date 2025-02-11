package common

// APIResponseSuccess constructs a typ.APIResponse of success.
func APIResponseSuccess(data interface{}) APIResponse {
	return APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// APIResponseError constructs a typ.APIResponse of error.
func APIResponseError(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
