package aws

import (
	"flag"
)

type AwsProviderConfiguration struct {
	IpProviderEnabled  *bool
	TagProviderEnabled *bool
}

func NewAwsProviderConfiguration() *AwsProviderConfiguration {
	config := AwsProviderConfiguration{}

	config.IpProviderEnabled = flag.Bool("provider-aws-ip-enabled", false, "Enables the aws provider for ip")
	config.TagProviderEnabled = flag.Bool("provider-aws-tags-enabled", false, "Enables the aws provider for tags ")

	return &config
}

func (this *AwsProviderConfiguration) Parse() {

}
