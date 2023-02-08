package kind

import (
	"crypto/rand"
	"fmt"
	"log"

	"sigs.k8s.io/kind/pkg/cluster"

	"github.com/aws-controllers-k8s/test-setup/pkg/config"
)

func generateRandomClusterName() *string {
	suffixLength := 5
	b := make([]byte, suffixLength)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	name := fmt.Sprintf("ack-test-%x", b)
	return &name
}

func CreateCluster(config *config.ClusterConfig, configDir string) error {
	createOptions := []cluster.CreateOption{}

	if config.KINDConfig != nil && config.KINDConfig.FileName != nil {
		createOptions = append(createOptions, cluster.CreateWithConfigFile(*config.KINDConfig.FileName))
	}

	name := config.Name
	if name == nil {
		name = generateRandomClusterName()
	}

	log.Printf("Creating KIND cluster with name \"%s\"\n", *name)

	options := cluster.ProviderWithDocker()

	provider := cluster.NewProvider(options, cluster.ProviderWithLogger(NewLogger()))
	if err := provider.Create(*name, createOptions...); err != nil {
		return err
	}

	return nil
}
