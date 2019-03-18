package event

type StartEvent struct {
	Id, parentId, Name, Address string
	Port                        int
	Tags                        map[string]string
}
