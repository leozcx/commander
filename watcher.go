package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"
	"github.com/denverdino/commander/registry"
	"github.com/samalba/dockerclient"
	"os"
	"os/signal"
	"syscall"
)

type EventContext struct {
	node       string
	client     *dockerclient.DockerClient
	etcdClient *etcd.Client
}

var EVENT_STATUS_MAP map[string]bool = map[string]bool{
	"create":  true,
	"destroy": true,
	"die":     true,
	"export":  true,
	"kill":    true,
	"pause":   true,
	"restart": true,
	"start":   true,
	"stop":    true,
	"unpause": true,
}

func eventCallback(e *dockerclient.Event, ec chan error, args ...interface{}) {
	context := args[0].(*EventContext)

	flag := EVENT_STATUS_MAP[e.Status]

	if flag {
		log.Info("Handling event:", e)
		containerId := e.Id
		registry.SetContainer(context.client, context.etcdClient, context.node, containerId)
	} else {
		log.Info("Ingore event:", e)
	}

}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	for _ = range sigChan {
		os.Exit(0)
	}
}

func watch(c *cli.Context) {
	tlsConfig, err := parseTlsConfig(c)

	if err != nil {
		log.Fatal(err)
	}

	// see https://github.com/codegangsta/cli/issues/160
	addr := c.String("addr")
	if addr == "" {
		log.Fatalf("addr required to access a docker daemon. See '%s watch --help'.", c.App.Name)
	}

	// see https://github.com/codegangsta/cli/issues/160
	node := c.String("node")
	if node == "" {
		log.Fatalf("node required to access a docker daemon. See '%s watch --help'.", c.App.Name)
	}

	// see https://github.com/codegangsta/cli/issues/160
	etcdURL := c.String("etcd")
	if etcdURL == "" {
		log.Fatalf("etcd required to access a etcd server. See '%s watch --help'.", c.App.Name)
	}

	docker, err := dockerclient.NewDockerClient(addr, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	etcdClient := etcd.NewClient([]string{etcdURL})

	//Prepare all the etcd entries for Docker Containers
	//TODO Sync the host status and etcd registry
	registry.SetHostWithContainers(docker, etcdClient, node)

	context := EventContext{
		node,
		docker,
		etcdClient,
	}

	//Start to monitoring
	docker.StartMonitorEvents(eventCallback, nil, &context)

	waitForInterrupt()
}
