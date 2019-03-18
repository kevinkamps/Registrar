package ifconfig

import (
	"flag"
)

type IfconfigProviderConfiguration struct {
	IpProviderEnabled *bool
}

func NewIfconfigProviderConfiguration() *IfconfigProviderConfiguration {
	config := IfconfigProviderConfiguration{}

	config.IpProviderEnabled = flag.Bool("provider-ifconfig-ip-enabled", false, "Enables the ifconfig.co provider for ip settings")

	return &config
}

func (this *IfconfigProviderConfiguration) Parse() {

}
