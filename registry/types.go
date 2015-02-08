package registry

import (
	"github.com/samalba/dockerclient"
	"time"
)

const INFO_POSTFIX = "/info"
const CONTAINER_PREFIX = "/registry/containers/"

type BaseResources struct {
	Kind string
	Id   string
	Name string
}

type ContainerResource struct {
	Kind        string
	Node        string
	Id          string
	Name        string
	ServiceName string
	Image       string
	Config      *dockerclient.ContainerConfig
	Created     string
	State       struct {
		Running    bool
		Paused     bool
		Restarting bool
		Pid        int
		ExitCode   int
		StartedAt  time.Time
		FinishedAt time.Time
		Ghost      bool
	}
	NetworkSettings struct {
		IpAddress   string
		IpPrefixLen int
		Gateway     string
		Bridge      string
		Ports       map[string][]dockerclient.PortBinding
	}
}

const HOST_PREFIX = "/registry/hosts/"

type HostResource struct {
	Kind string
	dockerclient.Info
}

const SERVICE_PREFIX = "/registry/services/"

type ServiceResource struct {
	Kind string
	dockerclient.Info
}
