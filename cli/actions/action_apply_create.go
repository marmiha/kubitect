package actions

import (
	"cli/tools/ansible"
	"cli/tools/terraform"
)

func create(c *Cluster) error {
	t, err := terraform.NewTerraform(c.Ctx, c.Path)

	if err != nil {
		return err
	}

	if err := t.Apply(); err != nil {
		return err
	}

	if err := c.Sync(); err != nil {
		return err
	}

	sshUser := string(*c.InfraConfig.Cluster.NodeTemplate.User)
	sshPKey := string(*c.InfraConfig.Cluster.NodeTemplate.SSH.PrivateKeyPath)

	k8sVersion := string(*c.NewConfig.Kubernetes.Version)

	if err := ansible.KubitectInit(c.Path, ansible.KUBESPRAY, ansible.GEN_NODES); err != nil {
		return err
	}

	if err := c.SetupKubesprayVE(); err != nil {
		return err
	}

	if err := ansible.HAProxy(c.Path, sshUser, sshPKey); err != nil {
		return err
	}

	if err := ansible.KubesprayCreate(c.Path, sshUser, sshPKey, k8sVersion); err != nil {
		return err
	}

	return ansible.KubitectFinalize(c.Path, sshUser, sshPKey)
}
