package types

import "github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"

type PoolLongDescription struct {
	CurrentCount  int  `json:"current_count"`
	Pool *pools.Pool   `json:"pool"`
	*APIResponse
}

func NewEmptyPoolLongDescription() *PoolLongDescription  {
	return &PoolLongDescription{}
}

func NewPoolLongDescription(version string) *PoolLongDescription {
	return &PoolLongDescription{APIResponse: NewApiResponse(version)}
}