package registrar

import (
	"kevinkamps/registrar/helper"
	"kevinkamps/registrar/provider"
	"kevinkamps/registrar/provider/aws"
	"kevinkamps/registrar/provider/ifconfig"
	"kevinkamps/registrar/provider/local"
	"kevinkamps/registrar/registrar/consul"
	"kevinkamps/registrar/registrar/event"
	"log"
	"sync"
)

type RegistrarService struct {
	registrars                   []Registrar
	ipProviders                  []provider.IpProvider
	tagProviders                 []provider.TagProvider
	networkProviderConfiguration *local.LocalNetworkProviderConfiguration
}

func NewRegistrarService(consulConfiguration *consul.Configuration,
	networkProviderConfiguration *local.LocalNetworkProviderConfiguration,
	ifconfigProviderConfiguration *ifconfig.IfconfigProviderConfiguration,
	awsProviderConfiguration *aws.AwsProviderConfiguration) *RegistrarService {
	service := RegistrarService{}

	service.networkProviderConfiguration = networkProviderConfiguration

	/**
	Registrars
	*/
	if *consulConfiguration.Enabled {
		log.Println("Registrar enabled: Consul")
		service.registrars = append(service.registrars, &consul.ConsulRegistrar{Configuration: consulConfiguration})
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

func (this *RegistrarService) Start() {
	var wg sync.WaitGroup
	for _, r := range this.registrars {
		wg.Add(1)
		go func(registrar Registrar) {
			registrar.Start()
			wg.Done()
		}(r)
	}
	wg.Wait()
}

func (this *RegistrarService) AddEvent(e event.Event) {
	for _, r := range this.registrars {
		if helper.IsInstanceOf(e, (*event.StartEvent)(nil)) {
			this.ProcessIpProviders(e.(*event.StartEvent))
			this.ProcessTagProviders(e.(*event.StartEvent))
		}
		r.AddEvent(e)
	}
}

func (this *RegistrarService) ProcessIpProviders(e *event.StartEvent) {
	for _, ipProvider := range this.ipProviders {
		ipProvider.AddAddress(e)
	}
}

func (this *RegistrarService) ProcessTagProviders(e *event.StartEvent) {
	for _, tagProvider := range this.tagProviders {
		tagProvider.AddTags(e)
	}
}
