package config

type Config struct {
	Cluster    *ClusterConfig `mapstructure:"cluster"`
	AWS        *AWSConfig     `mapstructure:"aws"`
	Tests      *TestConfig    `mapstructure:"tests"`
	Debug      *DebugConfig   `mapstructure:"debug"`
	LocalBuild *bool          `mapstructure:"local_build,omitempty"`
}

type ClusterConfig struct {
	Create        *bool              `mapstructure:"create"`
	Name          *string            `mapstructure:"name,omitempty"`
	K8sVersion    *string            `mapstructure:"k8s_version,omitempty"`
	Configuration *KINDConfiguration `mapstructure:"configuration"`
}

type KINDConfiguration struct {
	FileName              *string   `mapstructure:"file_name,omitempty"`
	AdditionalControllers []*string `mapstructure:"additional_controllers,omitempty"`
}

type AWSConfig struct {
	Profile        *string `mapstructure:"profile,omitempty"`
	TokenFile      *string `mapstructure:"token_file,omitempty"`
	Region         *string `mapstructure:"region,omitempty"`
	AssumedRoleARN *string `mapstructure:"assumed_role_arn"`
}

type TestConfig struct {
	RunLocally *bool     `mapstructure:"run_locally,omitempty"`
	Markers    []*string `mapstructure:"markers"`
	Methods    []*string `mapstructure:"methods"`
}

type DebugConfig struct {
	Enabled            *bool `mapstructure:"enabled,omitempty"`
	DumpControllerLogs *bool `mapstructure:"dump_controller_logs,omitempty"`
}
