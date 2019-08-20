package types

import (
	"fmt"
)

type ErrorMessage string

type APIResponse struct {
	ApiVersion string `json:"api_version"`
	Error ErrorMessage `json:"error"`
}

func NewApiResponse(version string) *APIResponse  {
	return &APIResponse{ApiVersion:version}
}

func NewFormattedErrorMsg(err error, format string, args ...interface{}) ErrorMessage {
	if err != nil {
		return ErrorMessage(fmt.Sprintf("%s : %s", fmt.Sprintf(format, args...), err.Error()))
	}
	return ErrorMessage(fmt.Sprintf(format, args))
}

func NewContextError(err error) ErrorMessage {
	return NewFormattedErrorMsg(err, "context not configured correctly")
}

func NewMissingParameterError(parameter string) ErrorMessage {
	return NewFormattedErrorMsg(fmt.Errorf("require parameter"), "missing parameter %s", parameter)
}
