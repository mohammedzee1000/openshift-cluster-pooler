package types

type StringResponse struct {
	Data  string `json:"data"`
	*APIResponse
}

func NewEmptyStringResponse() *StringResponse {
	return &StringResponse{}
}

func NewStringResponse(version string) *StringResponse {
	return &StringResponse{
		APIResponse: NewApiResponse(version),
	}
}