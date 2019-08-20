package types

import "github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"

type Clusters struct {
	Clusters *clusters.ClusterList
	*APIResponse
}

func NewEmptyClusters() *Clusters {
	return &Clusters{}
}

func NewClusters(version string) *Clusters {
	return &Clusters{
		APIResponse: NewApiResponse(version),
	}
}