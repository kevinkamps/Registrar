package registrar

import "kevinkamps/registrar/registrar/event"

type Registrar interface {
	Start()
	AddEvent(e event.Event)
}
