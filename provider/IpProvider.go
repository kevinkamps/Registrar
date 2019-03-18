package provider

import "kevinkamps/registrar/registrar/event"

type IpProvider interface {
	AddAddress(event *event.StartEvent)
}
