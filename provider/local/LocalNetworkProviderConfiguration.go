package local

import (
	"flag"
)

type LocalNetworkProviderConfiguration struct {
	IpProviderEnabled *bool
	UseIpv4           *bool
	UseIpv6           *bool
	InterfaceName     *string
}

func NewNetworkProviderConfiguration() *LocalNetworkProviderConfiguration {
	config := LocalNetworkProviderConfiguration{}

	config.IpProviderEnabled = flag.Bool("provider-local-network-ip-enabled", false, "Enables the network provider for ip settings")
	config.UseIpv4 = flag.Bool("provider-local-network-use-ipv4", true, "Has president over network-provider-use-ipv6")
	config.UseIpv6 = flag.Bool("provider-local-network-use-ipv6", false, "Can only be used if network-provider-use-ipv4 is set to false")
	config.InterfaceName = flag.String("provider-local-network-interface-name", "eth0", "Name of the interface to use")

	return &config
}

func (this *LocalNetworkProviderConfiguration) Parse() {

}
