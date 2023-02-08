package cmd

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

var (
	// Used for flags.
	cfgFile        string
	kindCfgDir     string
	kubeConfigPath string

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

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $GOPATH/src/github.com/aws-controllers-k8s/test-infra/test_config.yaml)")
	rootCmd.PersistentFlags().StringVar(&kindCfgDir, "kind-config-dir", fmt.Sprintf("%s/src/github.com/aws-controllers-k8s/test-infra/scripts/kind-configurations", goPath), "directory of the KIND configuration files (default is $GOPATH/src/github.com/aws-controllers-k8s/test-infra/scripts/kind-configurations)")

	if home := homedir.HomeDir(); home != "" {
		rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
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

	log.Println("Using config file: ", viper.ConfigFileUsed())
}
