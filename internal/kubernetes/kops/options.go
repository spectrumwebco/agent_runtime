package kops

import (
	"io"
)

type CreateClusterOptions struct {
	Name              string
	Zones             []string
	NodeCount         int
	MasterCount       int
	NodeSize          string
	MasterSize        string
	CloudProvider     string
	NetworkCIDR       string
	KubernetesVersion string
	SSHPublicKey      string
	DryRun            bool
	Output            string
	Yes               bool
	Stdout            io.Writer
	Stderr            io.Writer
}

type UpdateClusterOptions struct {
	Name             string
	Yes              bool
	CreateKubeconfig bool
	Output           string
	Stdout           io.Writer
	Stderr           io.Writer
}

type DeleteClusterOptions struct {
	Name   string
	Yes    bool
	Stdout io.Writer
	Stderr io.Writer
}

type ValidateClusterOptions struct {
	Name   string
	Output string
	Stdout io.Writer
	Stderr io.Writer
}

type GetClustersOptions struct {
	Output string
	Stdout io.Writer
	Stderr io.Writer
}

type ExportKubecfgOptions struct {
	Name       string
	KubeConfig string
	Admin      bool
	Stdout     io.Writer
	Stderr     io.Writer
}

type RollingUpdateOptions struct {
	Name               string
	Yes                bool
	Force              bool
	CloudOnly          bool
	MasterInterval     string
	NodeInterval       string
	InstanceGroupNames []string
	Stdout             io.Writer
	Stderr             io.Writer
}

type GetInstanceGroupsOptions struct {
	ClusterName string
	Output      string
	Stdout      io.Writer
	Stderr      io.Writer
}

type EditInstanceGroupOptions struct {
	ClusterName       string
	InstanceGroupName string
	Stdout            io.Writer
	Stderr            io.Writer
}

type GetSecretsOptions struct {
	Output string
	Stdout io.Writer
	Stderr io.Writer
}

type ToolboxDumpOptions struct {
	ClusterName string
	Output      string
	Stdout      io.Writer
	Stderr      io.Writer
}
