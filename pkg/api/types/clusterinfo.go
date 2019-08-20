package types

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
)

type ClusterInfo struct {
	Data    *clusters.Cluster
	*APIResponse
}

func NewEmptyClusterInfo() *ClusterInfo {
	return &ClusterInfo{}
}

func NewClusterInfo(version string) *ClusterInfo {
	return &ClusterInfo{
		APIResponse: NewApiResponse(version),
	}
}