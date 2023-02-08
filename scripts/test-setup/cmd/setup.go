package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws-controllers-k8s/test-setup/pkg/config"
	"github.com/aws-controllers-k8s/test-setup/pkg/kind"
	"github.com/aws-controllers-k8s/test-setup/pkg/kubeclient"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up the infrastructure for the tests",
	RunE:  setup,
}

func setup(cmd *cobra.Command, args []string) (err error) {
	config := &config.Config{}
	if err := viper.UnmarshalExact(config); err != nil {
		return fmt.Errorf("unable to decode config:\n%v", err)
	}

	if config.Cluster.Create != nil && *config.Cluster.Create {
		kind.CreateCluster(config.Cluster, kindCfgDir)
	} else {
		log.Println("Skipping create cluster")
	}

	client, err := kubeclient.NewClient(&kubeConfigPath)
	if err != nil {
		return err
	}

	if context, err := client.GetCurrentContext(); err != nil {
		return err
	} else {
		log.Println("Using context: ", context)
	}

	return nil
}
