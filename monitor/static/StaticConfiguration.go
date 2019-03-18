package static

import (
	"flag"
)

type Configuration struct {
	Enabled, LogChecksEnabled *bool
	ConfigPath                *string
	CheckDelay, CheckTimout   *int
}

func NewStaticConfiguration() *Configuration {
	config := Configuration{}

	config.Enabled = flag.Bool("monitor-static-enabled", false, "Enables the static monitor. Only read the config file at startup.")
	config.ConfigPath = flag.String("monitor-static-config-path", "./config.yml", "A path to the static config file")
	config.CheckDelay = flag.Int("monitor-static-check-delay", 10, "Checks every x seconds whether the service is reachable through the network (OSI layer 3 checks)")
	config.CheckTimout = flag.Int("monitor-static-check-timeout", 2, "Check timeout in seconds")
	config.LogChecksEnabled = flag.Bool("monitor-static-log-checks-enabled", false, "Logging layer 3 checks")

	return &config
}

func (this *Configuration) Parse() {

}
