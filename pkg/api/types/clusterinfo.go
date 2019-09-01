package types

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"time"
)

type ClusterInfo struct {
	Cluster    *clusters.Cluster `json:"Cluster"`
	ExpiresOn  time.Time  `json:"ExpiresOn"`
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