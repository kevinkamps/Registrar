package local

import (
	"kevinkamps/registrar/registry/event"
	"net"
)

type LocalNetworkProvider struct {
	ip            string
	Configuration *LocalNetworkProviderConfiguration
}

func NewLocalNetworkProvider(configuration *LocalNetworkProviderConfiguration) *LocalNetworkProvider {
	localNetworkProvider := LocalNetworkProvider{}

	localNetworkProvider.Configuration = configuration

	return &localNetworkProvider
}

func (this *LocalNetworkProvider) AddAddress(event *event.StartEvent) {

	if len(this.ip) == 0 {
		ifaces, _ := net.Interfaces()
		for _, i := range ifaces {
			if i.Name == *this.Configuration.InterfaceName {
				addrs, _ := i.Addrs()
				// handle err
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip == nil {
						continue
					}

					if *this.Configuration.UseIpv4 {
						this.ip = ip.To4().String()
						break
					}
					if *this.Configuration.UseIpv6 {
						this.ip = ip.To16().String()
						break
					}

				}
			}
		}
	}

	event.Address = this.ip
}
