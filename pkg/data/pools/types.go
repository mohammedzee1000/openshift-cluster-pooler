package pools

import "github.com/mohammedzee1000/openshift-cluster-pool/pkg/duration"

const (
	Pool_Key  = "Pool"
	Cleanup_Pool_Key = "Cleanup-Pool"
)

//Pool repersents generic of a pool
type Pool struct {
	Name                        string `json:"Name"`
	Description                 string `json:"Description"`
	Size                        int    `json:"Size"`
	MaxSize                     int    `json:"MaxSize"`
	UnusedClusterTimeout        duration.Duration    `json:"UnusedClusterTimeout"`
	UsedClusterTimeout          duration.Duration    `json:"UsedClusterTimeout"`
	ParallelProvisioning        int    `json:"ParallelProvisioning"`
	ParallelDeProvisioning      int    `json:"ParallelDeProvisioning"`
	ProvisionCommand            string `json:"ProvisionCommand"`
	DeProvisionCommand          string `json:"DeProvisionCommand"`
	ForceDeprovisionCommand     string `json:"ForceDeProvisionCommand"`
	ClusterURLCommand           string `json:"ClusterUrlCommand"`
	ClusterAdminUserCommand     string `json:"ClusterAdminUserCommand"`
	ClusterAdminPasswordCommand string `json:"ClusterAdminPasswordCommand"`
	ClusterCAFilePath           string `json:"ClusterCAFilePath"`
	ClusterCertFilePath         string `json:"ClusterCertFilePath"`
	ClusterKeyFilePath          string `json:"ClusterKeyFilePath"`
	ClusterExtraInfoCommand     string `json:"ClusterExtraInfoCommand"`
}

func NewEmptyPool() *Pool {
	return &Pool{}
}
