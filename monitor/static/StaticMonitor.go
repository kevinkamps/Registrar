package static

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"kevinkamps/registrar/registry"
	"kevinkamps/registrar/registry/event"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type StaticConfig struct {
	Applications []StaticApplication `yaml:"applications"`
}

type StaticApplication struct {
	Name     string            `yaml:"name"`
	Ip       string            `yaml:"ip"`
	Port     int               `yaml:"port"`
	Protocol string            `yaml:"protocol"`
	Tags     map[string]string `yaml:"tags"`
}

type Check struct {
	Ip           string
	EventId      string
	Application  StaticApplication
	deregistered bool
}

type StaticMonitor struct {
	RegistryService *registry.RegistryService
	Configuration   *Configuration
	checks          map[string]*Check
	hostname        string
}

func (this *StaticMonitor) getConf() *StaticConfig {
	config := StaticConfig{}

	yamlFile, err := ioutil.ReadFile(*this.Configuration.ConfigPath)
	if err != nil {
		log.Fatalf("Reading yaml file failed: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Parsing yaml file failed: %v", err)
	}

	for _, application := range config.Applications {
		if application.Protocol == "" {
			application.Protocol = "tcp"
		}
	}
	return &config
}

func (this *StaticMonitor) Start() {
	this.checks = make(map[string]*Check)

	config := this.getConf()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	this.hostname = hostname

	for _, application := range config.Applications {
		address := "127.0.0.1"
		if application.Ip != "" {
			address = application.Ip
		}

		e := this.registerApplication(address, application)
		this.RegistryService.ProcessIpProviders(e)
		this.checks[e.Id] = &Check{Application: application, Ip: e.Address, EventId: e.Id, deregistered: false}

	}

	for {
		for _, check := range this.checks {
			conn, err := net.DialTimeout(check.Application.Protocol, net.JoinHostPort(check.Ip, strconv.Itoa(check.Application.Port)), time.Duration(*this.Configuration.CheckTimout)*time.Second)

			if err != nil || conn == nil {
				if *this.Configuration.LogChecksEnabled {
					log.Println(fmt.Sprintf("Monitor - Static config: Check failed for %s %s:%s (timeout %s). Error: %s", check.Application.Protocol, check.Ip, strconv.Itoa(check.Application.Port), time.Duration(*this.Configuration.CheckTimout)*time.Second, err))
				}
				e := event.EndEvent{
					Id: check.EventId,
				}
				if !check.deregistered {
					this.RegistryService.AddEvent(&e)
					check.deregistered = true
				}
			} else {
				if *this.Configuration.LogChecksEnabled {
					log.Println(fmt.Sprintf("Monitor - Static config: Check succeeded for %s %s:%s (timeout %s)", check.Application.Protocol, check.Ip, strconv.Itoa(check.Application.Port), time.Duration(*this.Configuration.CheckTimout)*time.Second))
				}
				if check.deregistered {
					this.registerApplication(check.Ip, check.Application)
					check.deregistered = false
				}
			}

			if conn != nil {
				conn.Close()
			}
		}

		time.Sleep(time.Duration(*this.Configuration.CheckDelay) * time.Second)
	}
}

func (this *StaticMonitor) registerApplication(address string, application StaticApplication) *event.StartEvent {
	e := &event.StartEvent{
		Id:      fmt.Sprintf("registrar-static-%s-%s-%d", this.hostname, application.Name, application.Port),
		Name:    application.Name,
		Address: address,
		Port:    application.Port,
		Tags:    application.Tags,
	}

	this.RegistryService.AddEvent(e)

	return e
}
