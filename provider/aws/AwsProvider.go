package aws

import (
	"kevinkamps/registrar/provider"
	"kevinkamps/registrar/registry/event"
	"log"
)

type AwsProvider struct {
	localIpv4  *provider.WebRequestSingleValue
	hostname   *provider.WebRequestSingleValue
	instanceId *provider.WebRequestSingleValue
}

func NewAwsProvider() *AwsProvider {
	awsProvider := AwsProvider{}

	awsProvider.localIpv4 = &provider.WebRequestSingleValue{Url: "http://169.254.169.254/latest/meta-data/local-ipv4"}
	awsProvider.hostname = &provider.WebRequestSingleValue{Url: "http://169.254.169.254/latest/meta-data/hostname"}
	awsProvider.instanceId = &provider.WebRequestSingleValue{Url: "http://169.254.169.254/latest/meta-data/instance-id"}

	return &awsProvider
}

func (this *AwsProvider) AddAddress(event *event.StartEvent) {
	ip, err := this.localIpv4.GetValue()
	if err != nil {
		log.Printf("%s", err)
		return
	}
	event.Address = *ip
}

func (this *AwsProvider) AddTags(event *event.StartEvent) {
	instanceId, err := this.instanceId.GetValue()
	if err != nil {
		log.Printf("%s", err)
		return
	}

	hostname, err := this.hostname.GetValue()
	if err != nil {
		log.Printf("%s", err)
		return
	}

	event.Tags["InstanceId"] = *instanceId
	event.Tags["Hostname"] = *hostname
}
