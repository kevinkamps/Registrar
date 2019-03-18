# Registrar
Registrar automatically registers and deregisters applications or services to a registry. The only registry that is currently supported 
is a service discovery tool called [Consul](https://www.consul.io/). The Registrar can listen for docker containers starting and stopping as well as
reading from a static config file. It can use providers (AWS, WAN and local network) to resolve ips or provide additional tags when the application or service registrates.

## Project background
This project took a lot if inspiration from [Gliderlabs registrator](https://github.com/gliderlabs/registrator). 
I originally used it but found it to be lacking in resolving ips automatically in a cloud based environment. I could have added this functionality to that software but i
always wanted to learn GO and decided to give it a go (pun intended). My goal was to write something that i could use on my own personal servers 
as well as for a large customers i work for in a constantly changing high available AWS cloud environment. Also i wanted it to be more flexible than the original project.
I ended up with something that is easy to expand on and proved to be quite stable with a very low memory profile. So i decided to make it open source and share it
with all of you. 

## Features
The project consist out of three basic type of components. Monitors, Providers and Registries.

### Overview

*Monitors*
* Static config file
* [Docker](https://www.docker.com/)

*Providers*
* Local network ip
* WAN ip address
* [AWS](https://aws.amazon.com/) ip address or tags

*Registries*
* [Consul](https://www.consul.io/)


### Monitors
Monitors for applications. Technically you could use multiple monitors at the same time. Currently we support:
* Static config monitor
* Docker monitor

#### Static config
Monitors a static configuration file and checks if the applications is available and passing it to the registries.

The config file (by default `config.yml`) is using a yaml format. A example is provide below and should be self explanatory:
```yaml
applications:
- {
  name: "Application",
  port: 443,
  ip: "192.168.168.26",
  protocol: "tcp", # tcp|udp
  tags: {
    traefik-public.enable: "true",
    traefik-public.protocol: "https",
    traefik-public.frontend.entryPoints: "http,https",
    traefik-public.frontend.whiteList.sourceRange: "192.168.168.0/24,123.123.123.123/32"
  }
}
```

Please not that instead of setting the ip yourself it can be automatically resolved by any of the ip providers available.

Watching the config file for changes at runtime has not been implemented yet. So if you make changes to your config file please restart the Registrar. This feature will be added in the future.


#### Docker
Monitors containers on the local machine

Optional you can use the docker labels below on a container to influence the registration of your container. All labels
are prefixed with "REGISTRAR_" that would indicate that these labels are related to the registrar application.

Docker labels options:
 * `REGISTRAR_NAME`: "< name >". Name of the service to register
 * `REGISTRAR_INGORE`: "< true|false >". If set to true all services will be ignored.
 * `REGISTRAR_TAG_< name_of_the_tag >`: "< value of the tag >". Adds a tag can are used while registrating the application
 * `REGISTRAR_< private_port >_NAME`: "< name >". Name of the service to register
 * `REGISTRAR_< private_port >_IGNORE`: "< true|false >" If set to true this services will be ignored. If set to false it will not be ignored even if `REGISTRAR_INGORE` is set to true.
 * `REGISTRAR_< private_port >_TAG_< name_of_the_tag >`: "< value of the tag >". Adds a tag like `REGISTRAR_TAG_< name of the tag >` but will only add it the the specified port.



### Providers
Providers are data providers that can provide tags or ip addresses.

#### AWS provider (Ip and Tag provider)
Provides information about the [AWS](https://aws.amazon.com/) machine the registrar is running on. It can provide ip addresses and tags
* Ip: only local-ipv4 is supported right now
* tags: Would only add the following tags:
    * InstanceId
    * Hostname
    
#### Ifconfig provider (Ip provider)
Provides the WAN ip address as the ip address for the application

#### Local Network provider (Ip provider)
provides the ip address of a specified network interface as the applications ip address



### Registries
Registers the discovered applications to one of the implementations below.

#### consul
[Consul](https://www.consul.io/) Catalog (Registers to consul for service discovery usage) 


## CLI Options
```
Usage:
  -monitor-docker-api-version string
        Version of the api to use
  -monitor-docker-enabled
        Enables the docker monitor
  -monitor-docker-event-buffer-size int
        Max number of events to be buffered (default 1024)
  -monitor-docker-host string
        Docker host (default "unix:///var/run/docker.sock")
  -monitor-static-check-delay int
        Checks every x seconds whether the service is reachable through the network (OSI layer 3 checks) (default 10)
  -monitor-static-check-timeout int
        Check timeout in seconds (default 2)
  -monitor-static-config-path string
        A path to the static config file (default "./config.yml")
  -monitor-static-enabled
        Enables the static monitor. Only read the config file at startup.
  -monitor-static-log-checks-enabled
        Logging layer 3 checks
  -provider-aws-ip-enabled
        Enables the aws provider for ip
  -provider-aws-tags-enabled
        Enables the aws provider for tags 
  -provider-ifconfig-ip-enabled
        Enables the ifconfig.co provider for ip settings
  -provider-local-network-interface-name string
        Name of the interface to use (default "eth0")
  -provider-local-network-ip-enabled
        Enables the network provider for ip settings
  -provider-local-network-use-ipv4
        Has president over network-provider-use-ipv6 (default true)
  -provider-local-network-use-ipv6
        Can only be used if network-provider-use-ipv4 is set to false
  -registry-consul-check-deregister-after int
        deregister in seconds (default 60)
  -registry-consul-check-ttl int
        Ttl in seconds (default 10)
  -registry-consul-enabled
        Enables registration to consul
  -registry-consul-event-buffer-size int
        Max number of events to be buffered (default 1024)
  -registry-consul-log-ttl-passes-enabled
        Logging of ttl passes are enabled if set to true
  -registry-consul-url string
        Consul address (default "http://127.0.0.1:8500")
```

## Running as Docker container
### Commandline example (showing help)
```bash
docker run -d \
    --name=registrar \
    kevinkamps/registrar:latest \
      -help
```
### Commandline example (monitoring docker containers and registering them with consul. Registering them with the network interface ip):
```bash
docker run -d \
    --name=registrar \
    --net=host \
    --volume=/var/run/docker.sock:/tmp/docker.sock \
    kevinkamps/registrar:latest \
      -monitor-docker-enabled=true -provider-local-network-ip-enabled=true -provider-local-network-interface-name=eth0 -registries-consul-enabled=true -registries-consul-url=http://127.0.0.1:8500 
```

### Docker compose example:
```yaml
version: '2'
services:
  registrar:
    image: kevinkamps/registrar:latest
    command: -help
```

### Docker compose example (monitoring docker containers and registering them with consul. Registering them with the network interface ip):
```yaml
version: '2'
services:
  registrar:
    image: kevinkamps/registrar:latest
    network_mode: host
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: -monitor-docker-enabled=true -provider-local-network-ip-enabled=true -provider-local-network-interface-name=eth0 -registries-consul-enabled=true -registries-consul-url=http://127.0.0.1:8500
```

## Future ideas / plans
* Add notifiers after the registration process
* Add public ips to aws provider
* Add more registrars


## Development notes
GO mods must be enabled for this project before you can build it.
* Terminal: `export GO111MODULE=on`
* Intellij: Settings -> Language & Frameworks -> Go -> Go modules -> Enable Go Modules

## License

[GPL-3.0](https://choosealicense.com/licenses/gpl-3.0/)
