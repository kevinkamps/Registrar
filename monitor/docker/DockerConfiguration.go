package docker

import (
	"flag"
)

type Configuration struct {
	Enabled          *bool
	host             *string
	version          *string
	EventsBufferSize *int
}

func NewDockerConfiguration() *Configuration {
	config := Configuration{}

	config.Enabled = flag.Bool("monitor-docker-enabled", false, "Enables the docker monitor")
	config.host = flag.String("monitor-docker-host", `unix:///var/run/docker.sock`, "Docker host")
	config.version = flag.String("monitor-docker-api-version", ``, "Version of the api to use")
	config.EventsBufferSize = flag.Int("monitor-docker-event-buffer-size", 1024, "Max number of events to be buffered")

	return &config
}

func (this *Configuration) Parse() {

}
