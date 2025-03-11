package common

import "reflect"

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

// ToPtr converts a value of any type to a pointer of that type.
// It returns a pointer to the value if the value is not nil.
// If the value is nil, it returns a nil pointer.
func ToPtr[T any](v T) *T {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}
	return &v
}
