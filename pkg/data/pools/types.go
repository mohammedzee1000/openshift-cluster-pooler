package pools

import "github.com/mohammedzee1000/openshift-cluster-pool/pkg/duration"

const (
	Pool_Key  = "Pool"
	Cleanup_Pool_Key = "Cleanup-Pool"
)

//Pool repersents generic of a pool
type Pool struct {
	Name                        string `json:"name"`
	Description                 string `json:"description"`
	Size                        int    `json:"size"`
	MaxSize                     int    `json:"max-size"`
	UnusedClusterTimeout        duration.Duration    `json:"unused-cluster-timeout"`
	UsedClusterTimeout          duration.Duration    `json:"used-cluster-timeout"`
	ParallelProvisioning        int    `json:"parallel-provisioning"`
	ParallelDeProvisioning      int    `json:"parallel-de-provisioning"`
	ProvisionCommand            string `json:"provision-command"`
	DeProvisionCommand          string `json:"de-provision-command"`
	ForceDeprovisionCommand     string `json:"force-deprovision-command"`
	ClusterURLCommand           string `json:"cluster-url-command"`
	ClusterAdminUserCommand     string `json:"cluster-admin-user-command"`
	ClusterAdminPasswordCommand string `json:"cluster-admin-password-command"`
	ClusterCAFilePathString     string `json:"cluster-ca-file-path-string"`
	ClusterCertFilePathString   string `json:"cluster-cert-file-path-string"`
	ClusterKeyFilePathString    string `json:"cluster-key-file-path-string"`
	ClusterExtraInfoCommand     string `json:"cluster-extra-info-command"`
}

func NewEmptyPool() *Pool {
	return &Pool{}
}
