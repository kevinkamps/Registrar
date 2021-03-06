package docker

import (
	"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	"kevinkamps/registrar/registry"
	"kevinkamps/registrar/registry/event"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type DockerMonitor struct {
	RegistryService           *registry.RegistryService
	Configuration             *Configuration
	dockerApi                 *dockerapi.Client
	containerIdPrivatePortMap map[string][]int
	containerIdPublicPortMap  map[string][]int
}

func (this *DockerMonitor) registerAllCurrentRunningContainers() {
	containers, err := this.dockerApi.ListContainers(dockerapi.ListContainersOptions{All: true})
	assert(err)

	for _, container := range containers {
		containerInfo, err := this.dockerApi.InspectContainer(container.ID)
		assert(err)

		if containerInfo.State.Running {
			this.containerStarted(containerInfo)
		}
	}
}

func (this *DockerMonitor) containerStarted(container *dockerapi.Container) {
	log.Println(fmt.Sprintf("Monitor - Docker: Starting container detected. Container id: %s", container.ID))
	isNetworkHostMode := false
	if len(container.NetworkSettings.Networks) > 0 {
		if _, ok := container.NetworkSettings.Networks["host"]; ok {
			isNetworkHostMode = true
		}
	}

	if isNetworkHostMode {
		for portBinding := range container.Config.ExposedPorts {
			port, _ := strconv.Atoi(portBinding.Port())
			this.registerContainer(container, port, port)
		}
	} else {
		for privatePortProtocol, portBinding := range container.NetworkSettings.Ports {
			privatePort, _ := strconv.Atoi(privatePortProtocol.Port())

			for _, binding := range portBinding {
				publicPort, _ := strconv.Atoi(binding.HostPort)
				this.registerContainer(container, privatePort, publicPort)
			}
		}
	}
}

func (this *DockerMonitor) createEventId(containerId string, publicPort int, privatePort int) string {
	return fmt.Sprintf("registrar-docker-%s-%d:%d", containerId, publicPort, privatePort)
}

func (this *DockerMonitor) registerContainer(container *dockerapi.Container, privatePort int, publicPort int) {
	if !this.skipRegistration(container.Config.Labels, privatePort) {
		e := &event.StartEvent{
			Id:      this.createEventId(container.ID, publicPort, privatePort),
			Name:    this.getServiceName(container.Name, container.Config.Labels, privatePort, len(container.NetworkSettings.Ports) > 1),
			Address: "127.0.0.1",
			Port:    publicPort,
			Tags:    this.sanatizeLabels(container.Config.Labels, privatePort),
		}

		//Replacing address with first available docker address
		for _, network := range container.NetworkSettings.Networks {
			e.Address = network.IPAddress
			break
		}

		this.containerIdPrivatePortMap[container.ID] = append(this.containerIdPrivatePortMap[container.ID], privatePort)
		this.containerIdPublicPortMap[container.ID] = append(this.containerIdPublicPortMap[container.ID], publicPort)
		this.RegistryService.AddEvent(e)
	} else {
		log.Println(fmt.Sprintf("Monitor - Docker: Registration skipped because ignore flag was set. Container id: %s", container.ID))
	}
}
func (this *DockerMonitor) containerStopped(container *dockerapi.Container) {
	log.Println(fmt.Sprintf("Monitor - Docker: Stopping container detected. Container id: %s", container.ID))
	for i := range this.containerIdPrivatePortMap[container.ID] {
		privatePort := this.containerIdPrivatePortMap[container.ID][i]
		publicPort := this.containerIdPublicPortMap[container.ID][i]
		e := event.EndEvent{
			Id: this.createEventId(container.ID, publicPort, privatePort),
		}
		this.RegistryService.AddEvent(&e)
	}
	delete(this.containerIdPrivatePortMap, container.ID)
	delete(this.containerIdPublicPortMap, container.ID)
}

func (this *DockerMonitor) Start() {
	this.containerIdPrivatePortMap = make(map[string][]int)
	this.containerIdPublicPortMap = make(map[string][]int)

	var wg sync.WaitGroup
	os.Setenv("DOCKER_HOST", *this.Configuration.host)
	os.Setenv("DOCKER_API_VERSION", *this.Configuration.version)

	docker, err := dockerapi.NewClientFromEnv()
	if err != nil {
		assert(err)
	}
	this.dockerApi = docker

	events := make(chan *dockerapi.APIEvents, *this.Configuration.EventsBufferSize)
	assert(docker.AddEventListener(events))

	this.registerAllCurrentRunningContainers()

	wg.Add(1)
	go func() {
		for event := range events {
			switch event.Status {
			case "start":
				container, err := this.dockerApi.InspectContainer(event.ID)
				assert(err)
				this.containerStarted(container)
			case "die":
				container, err := this.dockerApi.InspectContainer(event.ID)
				assert(err)
				this.containerStopped(container)
			}
		}
	}()
	wg.Wait()
}

func (this *DockerMonitor) skipRegistration(labels map[string]string, port int) bool {
	if *this.Configuration.DefaultIgnoreEnabled {
		if value, ok := labels[fmt.Sprintf("REGISTRAR_%d_IGNORE", port)]; ok {
			if value == "false" {
				return false
			}
		}
		if value, ok := labels["REGISTRAR_IGNORE"]; ok {
			if value == "false" {
				return false
			}
		}
		return true
	} else {
		if value, ok := labels[fmt.Sprintf("REGISTRAR_%d_IGNORE", port)]; ok {
			return value == "true"
		}
		if value, ok := labels["REGISTRAR_IGNORE"]; ok {
			return value == "true"
		}
	}

	return false
}

func (this *DockerMonitor) getServiceName(name string, labels map[string]string, port int, registerWithPortSuffix bool) string {
	if value, ok := labels[fmt.Sprintf("REGISTRAR_%d_NAME", port)]; ok {
		return value
	}
	if value, ok := labels["REGISTRAR_NAME"]; ok {
		return value
	}
	name = name[1:len(name)]
	return fmt.Sprintf("%s-%d", name, port)
}

func (this *DockerMonitor) sanatizeLabels(labels map[string]string, port int) map[string]string {

	var labelsToKeep map[string]string = map[string]string{}

	for label, value := range labels {

		var prefix = "REGISTRAR_TAG_"
		if strings.HasPrefix(label, prefix) {
			tag := label[len(prefix):]
			labelsToKeep[tag] = value
			continue
		}

		prefix = fmt.Sprintf("REGISTRAR_%d_TAG_", port)
		if strings.HasPrefix(label, prefix) {
			tag := label[len(prefix):]
			labelsToKeep[tag] = value
			continue
		}
	}
	return labelsToKeep
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
