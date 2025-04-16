package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/kubernetes/kops"
	kopsModule "github.com/spectrumwebco/agent_runtime/pkg/modules/kops"
)

func NewKopsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kops",
		Short: "Kubernetes Operations (kOps) for cluster management",
		Long:  `kOps is the easiest way to get a production grade Kubernetes cluster up and running.`,
	}

	cmd.AddCommand(newKopsCreateClusterCommand())
	cmd.AddCommand(newKopsUpdateClusterCommand())
	cmd.AddCommand(newKopsDeleteClusterCommand())
	cmd.AddCommand(newKopsValidateClusterCommand())
	cmd.AddCommand(newKopsGetClustersCommand())
	cmd.AddCommand(newKopsExportKubecfgCommand())
	cmd.AddCommand(newKopsRollingUpdateCommand())
	cmd.AddCommand(newKopsGetInstanceGroupsCommand())
	cmd.AddCommand(newKopsEditInstanceGroupCommand())
	cmd.AddCommand(newKopsGetSecretsCommand())
	cmd.AddCommand(newKopsToolboxDumpCommand())
	cmd.AddCommand(newKopsInstallCommand())

	return cmd
}

func newKopsCreateClusterCommand() *cobra.Command {
	var (
		name              string
		zones             []string
		nodeCount         int
		masterCount       int
		nodeSize          string
		masterSize        string
		cloudProvider     string
		networkCIDR       string
		kubernetesVersion string
		sshPublicKey      string
		dryRun            bool
		output            string
		yes               bool
		stateStore        string
	)

	cmd := &cobra.Command{
		Use:   "create-cluster",
		Short: "Create a new Kubernetes cluster",
		Long:  `Create a new Kubernetes cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.CreateClusterOptions{
				Name:              name,
				Zones:             zones,
				NodeCount:         nodeCount,
				MasterCount:       masterCount,
				NodeSize:          nodeSize,
				MasterSize:        masterSize,
				CloudProvider:     cloudProvider,
				NetworkCIDR:       networkCIDR,
				KubernetesVersion: kubernetesVersion,
				SSHPublicKey:      sshPublicKey,
				DryRun:            dryRun,
				Output:            output,
				Yes:               yes,
				Stdout:            os.Stdout,
				Stderr:            os.Stderr,
			}

			return module.CreateCluster(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cluster")
	cmd.Flags().StringSliceVar(&zones, "zones", nil, "Zones in which to run the cluster")
	cmd.Flags().IntVar(&nodeCount, "node-count", 2, "Number of nodes")
	cmd.Flags().IntVar(&masterCount, "master-count", 1, "Number of masters")
	cmd.Flags().StringVar(&nodeSize, "node-size", "", "Size of the nodes")
	cmd.Flags().StringVar(&masterSize, "master-size", "", "Size of the masters")
	cmd.Flags().StringVar(&cloudProvider, "cloud", "", "Cloud provider (aws, gce, etc.)")
	cmd.Flags().StringVar(&networkCIDR, "network-cidr", "", "Network CIDR")
	cmd.Flags().StringVar(&kubernetesVersion, "kubernetes-version", "", "Kubernetes version")
	cmd.Flags().StringVar(&sshPublicKey, "ssh-public-key", "", "SSH public key")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Automatic yes to prompts")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("zones")

	return cmd
}

func newKopsUpdateClusterCommand() *cobra.Command {
	var (
		name             string
		yes              bool
		createKubeconfig bool
		output           string
		stateStore       string
	)

	cmd := &cobra.Command{
		Use:   "update-cluster",
		Short: "Update a Kubernetes cluster",
		Long:  `Update a Kubernetes cluster configuration using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.UpdateClusterOptions{
				Name:             name,
				Yes:              yes,
				CreateKubeconfig: createKubeconfig,
				Output:           output,
				Stdout:           os.Stdout,
				Stderr:           os.Stderr,
			}

			return module.UpdateCluster(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cluster")
	cmd.Flags().BoolVar(&yes, "yes", false, "Automatic yes to prompts")
	cmd.Flags().BoolVar(&createKubeconfig, "create-kube-config", false, "Create kubeconfig")
	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsDeleteClusterCommand() *cobra.Command {
	var (
		name       string
		yes        bool
		stateStore string
	)

	cmd := &cobra.Command{
		Use:   "delete-cluster",
		Short: "Delete a Kubernetes cluster",
		Long:  `Delete a Kubernetes cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.DeleteClusterOptions{
				Name:   name,
				Yes:    yes,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			return module.DeleteCluster(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cluster")
	cmd.Flags().BoolVar(&yes, "yes", false, "Automatic yes to prompts")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsValidateClusterCommand() *cobra.Command {
	var (
		name       string
		output     string
		stateStore string
	)

	cmd := &cobra.Command{
		Use:   "validate-cluster",
		Short: "Validate a Kubernetes cluster",
		Long:  `Validate a Kubernetes cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.ValidateClusterOptions{
				Name:   name,
				Output: output,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			return module.ValidateCluster(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cluster")
	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsGetClustersCommand() *cobra.Command {
	var (
		output     string
		stateStore string
	)

	cmd := &cobra.Command{
		Use:   "get-clusters",
		Short: "List all clusters",
		Long:  `List all clusters in the state store using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.GetClustersOptions{
				Output: output,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			return module.GetClusters(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	return cmd
}

func newKopsExportKubecfgCommand() *cobra.Command {
	var (
		name       string
		kubeconfig string
		admin      bool
		stateStore string
	)

	cmd := &cobra.Command{
		Use:   "export-kubecfg",
		Short: "Export kubeconfig",
		Long:  `Export kubeconfig for a cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.ExportKubecfgOptions{
				Name:       name,
				KubeConfig: kubeconfig,
				Admin:      admin,
				Stdout:     os.Stdout,
				Stderr:     os.Stderr,
			}

			return module.ExportKubecfg(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cluster")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	cmd.Flags().BoolVar(&admin, "admin", false, "Export admin credentials")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsRollingUpdateCommand() *cobra.Command {
	var (
		name               string
		yes                bool
		force              bool
		cloudOnly          bool
		masterInterval     string
		nodeInterval       string
		instanceGroupNames []string
		stateStore         string
	)

	cmd := &cobra.Command{
		Use:   "rolling-update",
		Short: "Perform a rolling update of a cluster",
		Long:  `Perform a rolling update of a cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.RollingUpdateOptions{
				Name:               name,
				Yes:                yes,
				Force:              force,
				CloudOnly:          cloudOnly,
				MasterInterval:     masterInterval,
				NodeInterval:       nodeInterval,
				InstanceGroupNames: instanceGroupNames,
				Stdout:             os.Stdout,
				Stderr:             os.Stderr,
			}

			return module.RollingUpdate(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cluster")
	cmd.Flags().BoolVar(&yes, "yes", false, "Automatic yes to prompts")
	cmd.Flags().BoolVar(&force, "force", false, "Force rolling update, even if no changes")
	cmd.Flags().BoolVar(&cloudOnly, "cloud-only", false, "Only update cloud resources")
	cmd.Flags().StringVar(&masterInterval, "master-interval", "", "Time to wait between restarting masters")
	cmd.Flags().StringVar(&nodeInterval, "node-interval", "", "Time to wait between restarting nodes")
	cmd.Flags().StringSliceVar(&instanceGroupNames, "instance-group", nil, "Instance groups to update")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsGetInstanceGroupsCommand() *cobra.Command {
	var (
		clusterName string
		output      string
		stateStore  string
	)

	cmd := &cobra.Command{
		Use:   "get-instance-groups",
		Short: "List all instance groups",
		Long:  `List all instance groups in a cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.GetInstanceGroupsOptions{
				ClusterName: clusterName,
				Output:      output,
				Stdout:      os.Stdout,
				Stderr:      os.Stderr,
			}

			return module.GetInstanceGroups(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&clusterName, "name", "", "Name of the cluster")
	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsEditInstanceGroupCommand() *cobra.Command {
	var (
		clusterName       string
		instanceGroupName string
		stateStore        string
	)

	cmd := &cobra.Command{
		Use:   "edit-instance-group",
		Short: "Edit an instance group",
		Long:  `Edit an instance group in a cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.EditInstanceGroupOptions{
				ClusterName:       clusterName,
				InstanceGroupName: instanceGroupName,
				Stdout:            os.Stdout,
				Stderr:            os.Stderr,
			}

			return module.EditInstanceGroup(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&clusterName, "name", "", "Name of the cluster")
	cmd.Flags().StringVar(&instanceGroupName, "instance-group", "", "Name of the instance group")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("instance-group")

	return cmd
}

func newKopsGetSecretsCommand() *cobra.Command {
	var (
		output     string
		stateStore string
	)

	cmd := &cobra.Command{
		Use:   "get-secrets",
		Short: "List all secrets",
		Long:  `List all secrets in a cluster using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.GetSecretsOptions{
				Output: output,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			return module.GetSecrets(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	return cmd
}

func newKopsToolboxDumpCommand() *cobra.Command {
	var (
		clusterName string
		output      string
		stateStore  string
	)

	cmd := &cobra.Command{
		Use:   "toolbox-dump",
		Short: "Dump cluster state",
		Long:  `Dump cluster state using kOps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			if stateStore != "" {
				module.SetStateStore(stateStore)
			}

			options := kops.ToolboxDumpOptions{
				ClusterName: clusterName,
				Output:      output,
				Stdout:      os.Stdout,
				Stderr:      os.Stderr,
			}

			return module.ToolboxDump(cmd.Context(), options)
		},
	}

	cmd.Flags().StringVar(&clusterName, "name", "", "Name of the cluster")
	cmd.Flags().StringVar(&output, "output", "", "Output format (yaml, json)")
	cmd.Flags().StringVar(&stateStore, "state", "", "State store location (e.g., s3://bucket-name)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newKopsInstallCommand() *cobra.Command {
	var (
		version     string
		installPath string
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install kOps binary",
		Long:  `Install kOps binary to the specified path.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module, err := kopsModule.NewModule(cfg)
			if err != nil {
				return err
			}

			return module.InstallBinary(cmd.Context(), version, installPath)
		},
	}

	cmd.Flags().StringVar(&version, "version", "latest", "Version of kOps to install")
	cmd.Flags().StringVar(&installPath, "path", "", "Path to install kOps binary")

	return cmd
}
