package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"kevinkamps/registrar/helper"
	"kevinkamps/registrar/registrar/event"
	"log"
	"sync"
	"time"
)

type ConsulRegistrar struct {
	events        chan event.Event
	Configuration *Configuration
	consulClient  *consulapi.Client
	registrations map[string]*consulapi.AgentServiceRegistration
}

func (this *ConsulRegistrar) initConsulConnection() {
	config := consulapi.DefaultConfig()

	//TODO fix tsl connection

	config.Address = this.Configuration.Url.Host
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul: ", this.Configuration.Url.Scheme)
	}

	this.consulClient = client
}

func (this *ConsulRegistrar) initTtlChecks(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		sleep := time.Duration(*this.Configuration.Ttl) * time.Second / 2

		log.Println(fmt.Sprintf("Consul: Sending ttl passes every %d nanoseconds is activated", sleep))
		for {
			for _, registration := range this.registrations {
				if *this.Configuration.LogTtlPassesEnabled {
					log.Println(fmt.Sprintf("Consul: Sending ttl pass for %s (%s)", registration.Check.Name, registration.Check.CheckID))
				}
				this.consulClient.Agent().PassTTL(registration.Check.CheckID, "Pass TTL")
			}
			time.Sleep(sleep)
		}
		wg.Done()
		log.Fatal("Consul: Sending ttl passes has stopped")
	}()
}

func (this *ConsulRegistrar) Start() {
	this.registrations = make(map[string]*consulapi.AgentServiceRegistration)

	var wg sync.WaitGroup

	this.initConsulConnection()
	this.initTtlChecks(&wg)

	// handle events
	this.events = make(chan event.Event, *this.Configuration.EventsBufferSize)
	for e := range this.events {

		if helper.IsInstanceOf(e, (*event.StartEvent)(nil)) {
			startEvent := e.(*event.StartEvent)
			log.Println("Consul: start event: ", startEvent)

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
			this.consulClient.Agent().ServiceRegister(registration)
		} else if helper.IsInstanceOf(e, (*event.EndEvent)(nil)) {
			endEvent := e.(*event.EndEvent)
			log.Println("Consul: end event: ", endEvent)
			delete(this.registrations, endEvent.Id)
			this.consulClient.Agent().ServiceDeregister(endEvent.Id)
		}
	}

	wg.Wait()
}

func (this *ConsulRegistrar) AddEvent(e event.Event) {
	this.events <- e
}
