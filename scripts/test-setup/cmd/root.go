package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "test-setup",
		Short: "A CLI tool for setting up and configuring the local test infrastructure",
		Long: `test-setup is a CLI tool that is used for setting up the local test infrastructer.
This includes creating new KIND clusters, creating the testing images, deploying
the controller to the cluster and running the tests.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $GOPATH/src/github.com/aws-controllers-k8s/test-infra/test_config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath("$GOPATH/src/github.com/aws-controllers-k8s/test-infra")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("test_config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
