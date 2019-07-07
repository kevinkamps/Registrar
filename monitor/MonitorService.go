package monitor

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"kevinkamps/registrar/monitor/docker"
	"kevinkamps/registrar/monitor/static"
	"kevinkamps/registrar/registry"
	"log"
	"sync"
)

type MonitorService struct {
	monitors []Monitor
}

func NewMonitorService(registryService *registry.RegistryService, dockerConfiguration *docker.Configuration, staticConfiguration *static.Configuration) *MonitorService {
	service := MonitorService{}

	/**
	Monitors
	*/
	if *dockerConfiguration.Enabled {
		log.Println(fmt.Sprintf("Monitor enabled: Docker. With config: %s", spew.Sdump(dockerConfiguration)))
		service.monitors = append(service.monitors, &docker.DockerMonitor{
			RegistryService: registryService,
			Configuration:   dockerConfiguration,
		})
	}
	if *staticConfiguration.Enabled {
		log.Println(fmt.Sprintf("Monitor enabled: Static. with config: %s", spew.Sdump(staticConfiguration)))
		service.monitors = append(service.monitors, &static.StaticMonitor{
			RegistryService: registryService,
			Configuration:   staticConfiguration,
		})
	}

	return &service
}

func (this *MonitorService) Start() {
	var wg sync.WaitGroup
	for _, m := range this.monitors {
		wg.Add(1)
		go func(monitor Monitor) {
			monitor.Start()
			wg.Done()
		}(m)
	}
	wg.Wait()
}
