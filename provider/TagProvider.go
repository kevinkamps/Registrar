package provider

import "kevinkamps/registrar/registry/event"

type TagProvider interface {
	AddTags(event *event.StartEvent)
}
