package pools

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
)

//deprovision destroys clusters resources using uuid and command
func (p Pool) deprovision(ctx *generic.Context, clusterid string, force bool) error {
	ctx.Log.Info("Pool deprovision", "deprovisioning clusters of pool %s", p.Name)
	var cmd string
	if !force {
		cmd = p.DeProvisionCommand
	} else {
		cmd = p.ForceDeprovisionCommand
	}
	c, err := clusters.ClusterByID(ctx, p.Name, clusterid)
	if err != nil {
		return err
	}
	c.State = clusters.ClusterDeProvisioning
	_ = c.Save(ctx)
	out, err := runCommand(clusterid, cmd)
	if err != nil {
		ctx.Log.Error("Pool deprovision", err,"failed to deprovision cluster %s, pool %s", c.ClusterID, p.Name)
		c.State = clusters.ClusterFailed
		_ = c.Save(ctx)
		PrintIfDebug(ctx.Debug, "deprovision command output", out)
		return err
	}
	c.Delete(ctx)
	PrintIfDebug(ctx.Debug, "deprovision command output", out)
	ctx.Log.Info("Pool deprovision", "successfully deprovisioned clusters %s, pool %s", c.ClusterID, p.Name)
	return nil
}
