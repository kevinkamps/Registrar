package registry

import (
	"kevinkamps/registrar/helper"
	"kevinkamps/registrar/provider"
	"kevinkamps/registrar/provider/aws"
	"kevinkamps/registrar/provider/ifconfig"
	"kevinkamps/registrar/provider/local"
	"kevinkamps/registrar/registry/consul"
	"kevinkamps/registrar/registry/event"
	"log"
	"sync"
)

type RegistryService struct {
	registries                   []Registry
	ipProviders                  []provider.IpProvider
	tagProviders                 []provider.TagProvider
	networkProviderConfiguration *local.LocalNetworkProviderConfiguration
}

func NewRegistryService(consulConfiguration *consul.Configuration,
	networkProviderConfiguration *local.LocalNetworkProviderConfiguration,
	ifconfigProviderConfiguration *ifconfig.IfconfigProviderConfiguration,
	awsProviderConfiguration *aws.AwsProviderConfiguration) *RegistryService {
	service := RegistryService{}

	service.networkProviderConfiguration = networkProviderConfiguration

	/**
	Registrars
	*/
	if *consulConfiguration.Enabled {
		log.Println("Registry enabled: Consul")
		service.registries = append(service.registries, &consul.ConsulRegistry{Configuration: consulConfiguration})
	}

	/**
	Ip providers
	*/
	if *networkProviderConfiguration.IpProviderEnabled {
		log.Println("IP Provider enabled: Local network")
		service.ipProviders = append(service.ipProviders, local.NewLocalNetworkProvider(networkProviderConfiguration))
	}
	if *ifconfigProviderConfiguration.IpProviderEnabled {
		log.Println("IP Provider enabled: Ifconfig.co")
		service.ipProviders = append(service.ipProviders, ifconfig.NewIfconfigProvider())
	}

	if *awsProviderConfiguration.IpProviderEnabled {
		log.Println("IP Provider enabled: AWS")
		service.ipProviders = append(service.ipProviders, aws.NewAwsProvider())
	}

	/**
	Tag providers
	*/
	if *awsProviderConfiguration.TagProviderEnabled {
		log.Println("Tag Provider enabled: AWS")
		service.tagProviders = append(service.tagProviders, aws.NewAwsProvider())
	}
	return &service
}

func (this *RegistryService) Start() {
	var wg sync.WaitGroup
	for _, r := range this.registries {
		wg.Add(1)
		go func(registry Registry) {
			registry.Start()
			wg.Done()
		}(r)
	}
	wg.Wait()
}

func (this *RegistryService) AddEvent(e event.Event) {
	for _, r := range this.registries {
		if helper.IsInstanceOf(e, (*event.StartEvent)(nil)) {
			this.ProcessIpProviders(e.(*event.StartEvent))
			this.ProcessTagProviders(e.(*event.StartEvent))
		}
		r.AddEvent(e)
	}
}

func (this *RegistryService) ProcessIpProviders(e *event.StartEvent) {
	for _, ipProvider := range this.ipProviders {
		ipProvider.AddAddress(e)
	}
}

func (this *RegistryService) ProcessTagProviders(e *event.StartEvent) {
	for _, tagProvider := range this.tagProviders {
		tagProvider.AddTags(e)
	}
}
