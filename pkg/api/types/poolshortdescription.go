package types

type PoolShortDescription struct {
	Description string `json:"description"`
	CurrentCount int `json:"current_count"`
	*APIResponse
}

func NewEmptyPoolShortDescription() *PoolShortDescription {
	return &PoolShortDescription{}
}

func NewPoolShortDescription(version string) *PoolShortDescription {
	return &PoolShortDescription{
		APIResponse: NewApiResponse(version),
	}
}