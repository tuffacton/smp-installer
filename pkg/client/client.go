package client

type Configuration struct {
	InfraConfig InfraConfiguration `yaml:"infra"`
}

type CloudProviderType string

const (
	AWSCloudProvider CloudProviderType = "aws"
	GCPCloudProvider CloudProviderType = "gcp"
)
