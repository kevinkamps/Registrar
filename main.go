package main

import (
	"flag"
	"fmt"
	"kevinkamps/registrar/configuration"
	"kevinkamps/registrar/monitor"
	"kevinkamps/registrar/monitor/docker"
	"kevinkamps/registrar/monitor/static"
	"kevinkamps/registrar/provider/aws"
	"kevinkamps/registrar/provider/ifconfig"
	"kevinkamps/registrar/provider/local"
	"kevinkamps/registrar/registrar"
	"kevinkamps/registrar/registrar/consul"
	"sync"
)

func main() {

	var configurations []configuration.Configuration

	/**
	Monitors
	*/
	dockerConfiguration := docker.NewDockerConfiguration()
	configurations = append(configurations, dockerConfiguration)

	staticConfiguration := static.NewStaticConfiguration()
	configurations = append(configurations, staticConfiguration)

	/**
	Providers
	*/
	networkProviderConfiguration := local.NewNetworkProviderConfiguration()
	configurations = append(configurations, networkProviderConfiguration)

	ifconfigProviderConfiguration := ifconfig.NewIfconfigProviderConfiguration()
	configurations = append(configurations, ifconfigProviderConfiguration)

	awsProviderConfiguration := aws.NewAwsProviderConfiguration()
	configurations = append(configurations, awsProviderConfiguration)

	/**
	registrar
	*/
	consulConfiguration := consul.NewConsulConfiguration()
	configurations = append(configurations, consulConfiguration)

	flag.Parse()
	for _, configuration := range configurations {
		configuration.Parse()
	}

	var wg sync.WaitGroup

	rs := registrar.NewRegistrarService(consulConfiguration, networkProviderConfiguration, ifconfigProviderConfiguration, awsProviderConfiguration)
	ms := monitor.NewMonitorService(rs, dockerConfiguration, staticConfiguration)

	wg.Add(1)
	go func() {
		rs.Start()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ms.Start()
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Shutting Down")
}
