package pools

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/logging"
	"strings"
	"time"
)

func (p Pool) gatherInfoOnSuccess(c *clusters.Cluster) (error) {
	var err error
	c.State = clusters.State_Success
	c.CreatedOn = time.Now()
	c.URL, err = p.getClusterURL(c.ClusterID)
	if err != nil {
		return err
	}
	c.AdminUser, err = p.getClusterAdminUser(c.ClusterID)
	if err != nil {
		return err
	}
	c.AdminPassword, err = p.getClusterAdminPassword(c.ClusterID)
	if err != nil {
		return err
	}
	c.CAFile, err = p.getClusterCAFile(c.ClusterID)
	if err != nil {
		if IsCommandMissingError(err) {
			c.CAFile = ""
			return nil
		}
		return err
	}
	c.CertFile, err = p.getClusterCertFile(c.ClusterID)
	if err != nil {
		if IsCommandMissingError(err) {
			c.CertFile = ""
			return nil
		}
		return err
	}
	c.KeyFile, err = p.getClusterKeyFile(c.ClusterID)
	if err != nil {
		if IsCommandMissingError(err) {
			c.KeyFile = ""
			return nil
		}
		return err
	}
	c.ExtraInfo, err = p.getClusterExtraInfo(c.ClusterID)
	if err != nil {
		if IsCommandMissingError(err) {
			c.ExtraInfo = ""
			return nil
		}
		return err
	}
	p.stripeol(c)
	return nil
}

func (p Pool) stripeol(c *clusters.Cluster)  {
	c.URL = strings.TrimRight(c.URL, "\n")
}

//provision provisions a clusters using provided command.
func (p Pool) provision(ctx *config.Context) error {
	logging.Info("Pool provision", fmt.Sprintf("provisioning clusters of pool %s", p.Name))
	clusterid := uuid.New().String()
	c := clusters.NewCluster(clusterid, p.Name)
	_ = c.Save(ctx)
	_, err := runCommand(clusterid, p.ProvisionCommand)
	if err != nil {
		c.State = clusters.State_Failed
		_ = c.Save(ctx)
		logging.Error("Pool provision", "failed provision of clusters of pool %s: %s", p.Name, err.Error())
		return err
	}
	err = p.gatherInfoOnSuccess(c)
	if err != nil {
		c.State = clusters.State_Failed
		_ = c.Save(ctx)
		logging.Error("Pool provision", "failed to provision cluster %s, pool %s", c.ClusterID, p.Name)
		return err
	}
	_ = c.Save(ctx)
	logging.Info("Pool provision", "successfully provisioned clusters %s, pool %s", c.ClusterID, p.Name)
	return nil
}

//deprovision destroys clusters resources using uuid and command
func (p Pool) deprovision(ctx *config.Context, clusterid string, force bool) error {
	logging.Info("Pool deprovision", "deprovisioning clusters of pool %s", p.Name)
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
	c.State = clusters.State_DeProvisioning
	_ = c.Save(ctx)
	_, err = runCommand(clusterid, cmd)
	if err != nil {
		logging.Error("Pool deprovision", "failed to deprovision cluster %s, pool %s", c.ClusterID, p.Name)
		c.State = clusters.State_Failed
		_ = c.Save(ctx)
		return err
	}
	c.Delete(ctx)
	logging.Info("Pool deprovision", "successfully deprovisioned clusters %s, pool %s", c.ClusterID, p.Name)
	return nil
}

//getClusterURL gets the url of specified clusters using provided uuid and command
func (p Pool) getClusterURL(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterURLCommand)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//getClusterAdminUser gets the admin user for the clusters using uuid and command
func (p Pool) getClusterAdminUser(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterAdminUserCommand)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//getClusterAdminPassword gets the password of the admin user using uuid and command
func (p Pool) getClusterAdminPassword(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterAdminPasswordCommand)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//getClusterCAFile gets the clusters CA file for clusters using uuid and command
func (p Pool) getClusterCAFile(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterCAFileCommand)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//getClusterCertFile gets the cert file of clusters using uuid and command
func (p Pool) getClusterCertFile(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterCertFileCommand)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//getClusterKeyFile gets key file of clusters using uuid and command
func (p Pool) getClusterKeyFile(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterKeyFileCommand)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//getClusterExtraInfo gets extra  custom info about clusters using uuid and command
func (p Pool) getClusterExtraInfo(clusterid string) (string, error) {
	out, err := runCommand(clusterid, p.ClusterExtraInfoCommand)
	if err != nil && !IsCommandMissingError(err) {
		return "", err
	}
	return string(out), nil
}