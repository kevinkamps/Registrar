package monitor

import (
	"kevinkamps/registrar/monitor/docker"
	"kevinkamps/registrar/monitor/static"
	"kevinkamps/registrar/registrar"
	"log"
	"sync"
)

type MonitorService struct {
	monitors []Monitor
}

func NewMonitorService(registrarService *registrar.RegistrarService, dockerConfiguration *docker.Configuration, staticConfiguration *static.Configuration) *MonitorService {
	service := MonitorService{}

	/**
	Monitors
	*/
	if *dockerConfiguration.Enabled {
		log.Println("Monitor enabled: Docker")
		service.monitors = append(service.monitors, &docker.DockerMonitor{
			RegistrarService: registrarService,
			Configuration:    dockerConfiguration,
		})
	}
	if *staticConfiguration.Enabled {
		log.Println("Monitor enabled: Static")
		service.monitors = append(service.monitors, &static.StaticMonitor{
			RegistrarService: registrarService,
			Configuration:    staticConfiguration,
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
