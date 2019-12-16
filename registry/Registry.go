package registry

import "kevinkamps/registrar/registry/event"

type Registry interface {
	Start()
	Init()
	AddEvent(e event.Event)
}
