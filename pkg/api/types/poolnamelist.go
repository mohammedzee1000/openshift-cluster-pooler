package types

type PoolNameList struct {
	Data   []string `json:"data"`
	*APIResponse
}

func NewEmptyPoolNameList() *PoolNameList {
	return &PoolNameList{}
}

func NewPoolNameList(version string) *PoolNameList {
	return &PoolNameList{
		APIResponse: NewApiResponse(version),
	}
}