package docker

import (
	"flag"
)

type Configuration struct {
	Enabled              *bool
	DefaultIgnoreEnabled *bool
	host                 *string
	version              *string
	EventsBufferSize     *int
}

func NewDockerConfiguration() *Configuration {
	config := Configuration{}

	config.Enabled = flag.Bool("monitor-docker-enabled", false, "Enables the docker monitor")
	config.DefaultIgnoreEnabled = flag.Bool("monitor-docker-default-ignore", false, "Ignores by default everything unless REGISTRATOR_IGNORE or REGISTRATOR_<port>_IGNORE is explicitly set to false. This inverts the default behaviour, which is registrating everything unless the ignore flags are set to true")
	config.host = flag.String("monitor-docker-host", `unix:///var/run/docker.sock`, "Docker host")
	config.version = flag.String("monitor-docker-api-version", ``, "Version of the api to use")
	config.EventsBufferSize = flag.Int("monitor-docker-event-buffer-size", 1024, "Max number of events to be buffered")

	return &config
}

func (this *Configuration) Parse() {

}
