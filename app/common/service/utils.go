package service

import "github.com/tundrawork/stargate/app/common/typ"

// APIResponseSuccess constructs a typ.APIResponse of success.
func APIResponseSuccess(data interface{}) typ.APIResponse {
	return typ.APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// APIResponseError constructs a typ.APIResponse of error.
func APIResponseError(code int, message string) typ.APIResponse {
	return typ.APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
