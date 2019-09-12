package clusters

type ClusterState string

const (
	ClusterProvisioning   ClusterState = "Provisioning"
	ClusterFailed         ClusterState = "Failed"
	ClusterSuccess        ClusterState = "Success"
	ClusterUsed           ClusterState = "Used"
	ClusterDeProvisioning ClusterState = "Deprovisioning"
	ClusterCleanup        ClusterState = "Cleanup"
	ClusterReturned       ClusterState = "Returned"
)

func (cs ClusterState) String() string {
	return  string(cs)
}
