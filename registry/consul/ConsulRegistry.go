package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"kevinkamps/registrar/helper"
	"kevinkamps/registrar/registry/event"
	"log"
	"sync"
	"time"
)

type ConsulRegistry struct {
	events        chan event.Event
	Configuration *Configuration
	consulClient  *consulapi.Client
	registrations map[string]*consulapi.AgentServiceRegistration
}

func (this *ConsulRegistry) initConsulConnection() {
	config := consulapi.DefaultConfig()

	//TODO fix tsl connection

	config.Datacenter = *this.Configuration.Datacenter
	config.Address = this.Configuration.Url.Host
	config.Scheme = this.Configuration.Url.Scheme
	log.Println(fmt.Sprintf("Registry - Consul: Connecting to: %s ", this.Configuration.Url))
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("Registry - Consul: ", this.Configuration.Url.Scheme)
	}

	this.consulClient = client
}

func (this *ConsulRegistry) initTtlChecks(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		sleep := time.Duration(*this.Configuration.Ttl) * time.Second / 2

		log.Println(fmt.Sprintf("Registry - Consul: Sending ttl passes every %d nanoseconds is activated", sleep))
		for {
			for _, registration := range this.registrations {
				if *this.Configuration.LogTtlPassesEnabled {
					log.Println(fmt.Sprintf("Registry - Consul: Sending ttl pass for %s (%s)", registration.Check.Name, registration.Check.CheckID))
				}
				this.consulClient.Agent().PassTTL(registration.Check.CheckID, "Pass TTL")
			}
			time.Sleep(sleep)
		}
		wg.Done()
		log.Fatal("Registry - Consul: Sending ttl passes has stopped")
	}()
}

func (this *ConsulRegistry) Start() {
	this.registrations = make(map[string]*consulapi.AgentServiceRegistration)

	var wg sync.WaitGroup

	this.initConsulConnection()
	this.initTtlChecks(&wg)

	// handle events
	this.events = make(chan event.Event, *this.Configuration.EventsBufferSize)
	for e := range this.events {

		if helper.IsInstanceOf(e, (*event.StartEvent)(nil)) {
			startEvent := e.(*event.StartEvent)
			log.Println("Registry - Consul: start event: ", startEvent)

			registration := new(consulapi.AgentServiceRegistration)
			registration.ID = startEvent.Id
			registration.Name = startEvent.Name
			registration.Address = startEvent.Address
			registration.Port = startEvent.Port
			for k, v := range startEvent.Tags {
				registration.Tags = append(registration.Tags, k+"="+v)
			}

			registration.Check = new(consulapi.AgentServiceCheck)
			registration.Check.CheckID = startEvent.Id
			registration.Check.Name = startEvent.Name
			registration.Check.TTL = fmt.Sprintf("%ds", *this.Configuration.Ttl)
			registration.Check.DeregisterCriticalServiceAfter = fmt.Sprintf("%ds", *this.Configuration.DeregisterCriticalServiceAfter)

			this.registrations[registration.ID] = registration
			err := this.consulClient.Agent().ServiceRegister(registration)
			if err != nil {
				log.Println("Registry - Consul: Error registering service: ", err)
			}

		} else if helper.IsInstanceOf(e, (*event.EndEvent)(nil)) {
			endEvent := e.(*event.EndEvent)
			log.Println("Registry - Consul: end event: ", endEvent)
			delete(this.registrations, endEvent.Id)
			err := this.consulClient.Agent().ServiceDeregister(endEvent.Id)
			if err != nil {
				log.Println("Registry - Consul: Error registering service: ", err)
			}
		}
	}

	wg.Wait()
}

func (this *ConsulRegistry) AddEvent(e event.Event) {
	this.events <- e
}
