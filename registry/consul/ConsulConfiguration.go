package consul

import (
	"flag"
	"log"
	"net/url"
)

type Configuration struct {
	url                            *string
	Url                            *url.URL
	Enabled, LogTtlPassesEnabled   *bool
	Ttl                            *int
	DeregisterCriticalServiceAfter *int
	EventsBufferSize               *int
	Datacenter                     *string
	Token                          *string
}

func NewConsulConfiguration() *Configuration {
	config := Configuration{}

	config.Enabled = flag.Bool("registry-consul-enabled", false, "Enables registration to consul")
	config.url = flag.String("registry-consul-url", "http://127.0.0.1:8500", "Consul address")
	config.Ttl = flag.Int("registry-consul-check-ttl", 10, "Ttl in seconds")
	config.DeregisterCriticalServiceAfter = flag.Int("registry-consul-check-deregister-after", 60, "deregister in seconds")
	config.EventsBufferSize = flag.Int("registry-consul-event-buffer-size", 1024, "Max number of events to be buffered")
	config.LogTtlPassesEnabled = flag.Bool("registry-consul-log-ttl-passes-enabled", false, "Logging of ttl passes are enabled if set to true")
	config.Datacenter = flag.String("registry-consul-datacenter", "dc1", "Consul datacenter")
	config.Token = flag.String("registry-consul-token", "", "Token is used to provide a per-request ACL token")

	return &config
}

func (this *Configuration) Parse() {
	parsedUrl, err := this.Url.Parse(*this.url)
	assert(err)
	this.Url = parsedUrl
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
