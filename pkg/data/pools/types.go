package pools

const Pool_Key  = "Pool"

//Pool repersents config of a pool
type Pool struct {
	Name                        string `json:"name"`
	Description                 string `json:"description"`
	Size                        int    `json:"size"`
	MaxSize                     int    `json:"max-size"`
	UnusedClusterTimeout        int    `json:"unused-cluster-timeout"`
	UsedClusterTimeout          int    `json:"used-cluster-timeout"`
	ParallelProvisioning        int    `json:"parallel-provisioning"`
	ParallelDeProvisioning      int    `json:"parallel-de-provisioning"`
	ProvisionCommand            string `json:"provision-command"`
	DeProvisionCommand          string `json:"de-provision-command"`
	ForceDeprovisionCommand     string `json:"force-deprovision-command"`
	ClusterURLCommand           string `json:"cluster-url-command"`
	ClusterAdminUserCommand     string `json:"cluster-admin-user-command"`
	ClusterAdminPasswordCommand string `json:"cluster-admin-password-command"`
	ClusterCAFileCommand        string `json:"cluster-ca-file-command"`
	ClusterCertFileCommand      string `json:"cluster-cert-file-command"`
	ClusterKeyFileCommand       string `json:"cluster-key-file-command"`
	ClusterExtraInfoCommand     string `json:"cluster-extra-info-command"`
}

func NewEmptyPool() *Pool {
	return &Pool{}
}
