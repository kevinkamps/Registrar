package provider

import "kevinkamps/registrar/registrar/event"

type TagProvider interface {
	AddTags(event *event.StartEvent)
}
