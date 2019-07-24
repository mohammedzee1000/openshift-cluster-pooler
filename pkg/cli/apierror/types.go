package apierror

type APIError struct {
	Errcode  int `json:"errcode"`
	Message   string `json:"message"`
}

func NewNotFoundError(msg string) *APIError {
	return &APIError{1, msg}
}

func IsNotFoundError(err *APIError) bool {
	if err.Errcode == 1 {
		return true
	}
	return false
}

func NewContextError() *APIError {
	return &APIError{2, "context not configured correctly"}
}

func IsContextError(err *APIError) bool {
	if err.Errcode == 2 {
		return true
	}
	return false
}

func NewListError(msg string) *APIError {
	return &APIError{3, msg}
}

func IsListError(err *APIError) bool {
	if err.Errcode == 3 {
		return true
	}
	return false
}