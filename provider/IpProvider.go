package provider

import "kevinkamps/registrar/registry/event"

type IpProvider interface {
	AddAddress(event *event.StartEvent)
}
