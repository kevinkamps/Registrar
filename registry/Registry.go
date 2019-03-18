package registry

import "kevinkamps/registrar/registry/event"

type Registry interface {
	Start()
	AddEvent(e event.Event)
}
