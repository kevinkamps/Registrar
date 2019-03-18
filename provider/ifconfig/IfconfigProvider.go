package ifconfig

import (
	"kevinkamps/registrar/provider"
	"kevinkamps/registrar/registrar/event"
	"log"
	"strings"
)

type IfconfigProvider struct {
	ip            *provider.WebRequestSingleValue
	Configuration *IfconfigProviderConfiguration
}

func NewIfconfigProvider() *IfconfigProvider {
	ifconfigProvider := IfconfigProvider{}

	ifconfigProvider.ip = &provider.WebRequestSingleValue{Url: "https://ifconfig.co/ip"}

	return &ifconfigProvider
}

func (this *IfconfigProvider) AddAddress(event *event.StartEvent) {
	ip, err := this.ip.GetValue()
	if err != nil {
		log.Printf("%s", err)
		return
	}
	event.Address = strings.Replace(*ip, "\n", "", -1)
}
