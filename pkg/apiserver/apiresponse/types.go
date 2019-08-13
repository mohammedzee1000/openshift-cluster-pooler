package apiresponse

type APIError struct {
	Errcode  int `json:"errcode"`
	Message   string `json:"message"`
}

type APIResponse struct {
	Error *APIError `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

func NewApiResponse() *APIResponse  {
	return &APIResponse{}
}