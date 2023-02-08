package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws-controllers-k8s/test-setup/pkg/config"
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

	return nil
}
